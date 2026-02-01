package main

import (
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/internal/utils/help"
	hls_filesharer_client "github.com/number571/hidden-lake/pkg/api/services/filesharer/client"
	"github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
)

var (
	gFlags = flag.NewFlagsBuilder(
		flag.NewFlagBuilder("-v", "--version").
			WithDescription("print version of application"),
		flag.NewFlagBuilder("-h", "--help").
			WithDescription("print information about application"),
		flag.NewFlagBuilder("-p", "--path").
			WithDescription("set path to load files").
			WithDefinedValue("."),
		flag.NewFlagBuilder("-s", "--service").
			WithDescription("set internal address of the HLS").
			WithDefinedValue("localhost:9541"),
		flag.NewFlagBuilder("-t", "--type").
			WithDescription("set type [local|personal|public]").
			WithDefinedValue("public"),
		flag.NewFlagBuilder("-f", "--friend").
			WithDescription("set alias name of the friend").
			WithDefinedValue(""),
		flag.NewFlagBuilder("-d", "--do").
			WithDescription("set runner [list|info|download|upload|delete]").
			WithDefinedValue(""),
		flag.NewFlagBuilder("-a", "--arg").
			WithDescription("set argument for runner <page|file>").
			WithDefinedValue(""),
	).Build()
)

func main() {
	args := os.Args[1:]
	if ok := gFlags.Validate(args); !ok {
		fmt.Println("args invalid")
		os.Exit(1)
	}

	if gFlags.Get("-v").GetBoolValue(args) {
		fmt.Println(build.GetVersion())
		return
	}

	if gFlags.Get("-h").GetBoolValue(args) {
		help.Println(settings.CAppFullName+"-cli", settings.CAppDescription, gFlags)
		return
	}

	ctx := context.Background()
	if err := runFunction(ctx, args); err != nil {
		fmt.Println("\n", err)
		os.Exit(1)
	}
}

func runFunction(pCtx context.Context, pArgs []string) error {
	inputPath := strings.TrimSuffix(gFlags.Get("-p").GetStringValue(pArgs), "/")
	if err := os.MkdirAll(inputPath, 0700); err != nil {
		return errors.Join(ErrMkdirPath, err)
	}

	hlfClient := hls_filesharer_client.NewClient(
		hls_filesharer_client.NewRequester(
			gFlags.Get("-s").GetStringValue(pArgs),
			&http.Client{Timeout: build.GetSettings().GetHttpCallbackTimeout()},
		),
	)

	stgType := gFlags.Get("-t").GetStringValue(pArgs)
	friend := gFlags.Get("-f").GetStringValue(pArgs)
	do := gFlags.Get("-d").GetStringValue(pArgs)

	var (
		isLocal    = false
		isPersonal = false
	)
	switch stgType {
	case "local":
		isLocal = true
	case "personal":
		isPersonal = true
	case "public":
		// used by default
	default:
		fmt.Println("AAA", stgType)
		return ErrUnknownStorageType
	}

	switch do {
	case "list":
		var (
			fileInfoList dto.IFileInfoList
			err          error
		)
		page := gFlags.Get("-a").GetInt64Value(pArgs)
		if isLocal {
			fileInfoList, err = hlfClient.GetLocalList(pCtx, friend, uint64(page)) // nolint:gosec
		} else {
			fileInfoList, err = hlfClient.GetRemoteList(pCtx, friend, uint64(page), isPersonal) // nolint:gosec
		}
		if err != nil {
			return err
		}
		printFileInfoList(fileInfoList)
	case "info":
		var (
			fileInfo dto.IFileInfo
			err      error
		)
		fileName := gFlags.Get("-a").GetStringValue(pArgs)
		if isLocal {
			fileInfo, err = hlfClient.GetLocalFileInfo(pCtx, friend, fileName)
		} else {
			fileInfo, err = hlfClient.GetRemoteFileInfo(pCtx, friend, fileName, isPersonal)
		}
		if err != nil {
			return err
		}
		printFileInfo(fileInfo)
	case "download":
		var (
			tmpFile *os.File
			err     error
		)

		fileName := gFlags.Get("-a").GetStringValue(pArgs)
		tmpFile, err = os.CreateTemp("", fmt.Sprintf("%s-*.tmp", fileName)) // nolint: perfsprint
		if err != nil {
			return err
		}

		defer func() {
			_ = tmpFile.Close()
			_ = os.Remove(tmpFile.Name())
		}()

		if isLocal { // nolint: nestif
			err := hlfClient.GetLocalFile(
				&processWriter{fW: tmpFile},
				pCtx,
				friend,
				fileName,
			)
			if err != nil {
				return err
			}
		} else {
			inProcess, recvFileHash, err := hlfClient.GetRemoteFile(
				&processWriter{fW: tmpFile},
				pCtx,
				friend,
				fileName,
				isPersonal,
			)
			if err != nil {
				return err
			}
			if inProcess {
				fmt.Println("\nprocessing...")
				return nil
			}

			gotFileHash, err := getFileHash(tmpFile.Name())
			if err != nil {
				return err
			}
			if recvFileHash != gotFileHash {
				return ErrHashIsInvalid
			}
		}

		fullPath := filepath.Join(inputPath, fileName)
		if err := copyFile(fullPath, tmpFile); err != nil {
			return err
		}

		fmt.Println("\ndone!")
	case "upload":
		if !isLocal {
			return ErrAvailableOnlyForTypeLocal
		}

		fileName := gFlags.Get("-a").GetStringValue(pArgs)
		fullPath := filepath.Join(inputPath, fileName)

		file, err := os.Open(fullPath) // nolint: gosec
		if err != nil {
			return err
		}
		defer func() { _ = file.Close() }()

		if err := hlfClient.PutLocalFile(pCtx, friend, fileName, file); err != nil {
			return err
		}
	case "delete":
		if !isLocal {
			return ErrAvailableOnlyForTypeLocal
		}

		fileName := gFlags.Get("-a").GetStringValue(pArgs)
		if err := hlfClient.DelLocalFile(pCtx, friend, fileName); err != nil {
			return err
		}
	default:
		return errors.Join(ErrUnknownAction, errors.New(do)) // nolint:err113
	}

	return nil
}

type processWriter struct {
	fP uint64
	fW io.Writer
}

func (p *processWriter) Write(b []byte) (n int, err error) {
	n, err = p.fW.Write(b)
	p.fP += uint64(n) // nolint: gosec
	fmt.Printf("\rdownloading ... %dB", p.fP)
	return n, err
}

func printFileInfoList(pFileInfoList dto.IFileInfoList) {
	list := pFileInfoList.GetList()
	for _, info := range list {
		printFileInfo(info)
	}
}

func printFileInfo(pFileInfo dto.IFileInfo) {
	fmt.Printf("Name: %s\nHash: %s\nSize: %d\n\n", pFileInfo.GetName(), pFileInfo.GetHash(), pFileInfo.GetSize())
}

func getFileHash(filename string) (string, error) {
	f, err := os.Open(filename) //nolint:gosec
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()
	h := sha512.New384()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return encoding.HexEncode(h.Sum(nil)), nil
}

func copyFile(dstFilePath string, tmpFile *os.File) error {
	dstFile, err := os.OpenFile(dstFilePath, os.O_CREATE|os.O_WRONLY, 0600) // nolint: gosec
	if err != nil {
		return err
	}
	defer func() { _ = dstFile.Close() }()
	if _, err := tmpFile.Seek(0, io.SeekStart); err != nil {
		return err
	}
	_, err = io.Copy(dstFile, tmpFile)
	return err
}
