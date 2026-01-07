package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/number571/hidden-lake/build"
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	hls_filesharer_client "github.com/number571/hidden-lake/internal/services/filesharer/pkg/client"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/stream"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/internal/utils/help"
)

var (
	gFlags = flag.NewFlagsBuilder(
		flag.NewFlagBuilder("-v", "--version").
			WithDescription("print version of application"),
		flag.NewFlagBuilder("-h", "--help").
			WithDescription("print information about application"),
		flag.NewFlagBuilder("-p", "--path").
			WithDescription("set path to config, database files").
			WithDefinedValue("."),
		flag.NewFlagBuilder("-r", "--retry").
			WithDescription("retry number on load chunk of file").
			WithDefinedValue("3"),
		flag.NewFlagBuilder("-k", "--kernel").
			WithDescription("set internal address of the HLK").
			WithDefinedValue("localhost:9572"),
		flag.NewFlagBuilder("-f", "--friend").
			WithDescription("set alias name of the friend").
			WithDefinedValue(""),
		flag.NewFlagBuilder("-d", "--do").
			WithDescription("set runner [list|info|load]").
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
		help.Println(settings.CAppFullName, settings.CAppDescription, gFlags)
		return
	}

	ctx := context.Background()
	if err := runFunction(ctx, args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runFunction(pCtx context.Context, pArgs []string) error {
	inputPath := strings.TrimSuffix(gFlags.Get("-p").GetStringValue(pArgs), "/")
	if err := os.MkdirAll(inputPath, 0700); err != nil {
		return errors.Join(ErrMkdirPath, err)
	}

	retryNum := gFlags.Get("-r").GetInt64Value(pArgs)
	if retryNum < 0 {
		return ErrRetryNum
	}

	hlkClient := hlk_client.NewClient(
		hlk_client.NewBuilder(),
		hlk_client.NewRequester(
			gFlags.Get("-k").GetStringValue(pArgs),
			&http.Client{Timeout: build.GetSettings().GetHttpCallbackTimeout()},
		),
	)

	hlfClient := hls_filesharer_client.NewClient(
		hls_filesharer_client.NewBuilder(),
		hls_filesharer_client.NewRequester(hlkClient),
	)

	friend := gFlags.Get("-f").GetStringValue(pArgs)
	do := gFlags.Get("-d").GetStringValue(pArgs)

	switch do {
	case "list":
		page := gFlags.Get("-a").GetInt64Value(pArgs)
		fileInfoList, err := hlfClient.GetListFiles(pCtx, friend, uint64(page)) // nolint:gosec
		if err != nil {
			return err
		}
		fileInfoListStr, err := json.MarshalIndent(fileInfoList, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(fileInfoListStr))
	case "info":
		fileName := gFlags.Get("-a").GetStringValue(pArgs)
		fileInfo, err := hlfClient.GetFileInfo(pCtx, friend, fileName)
		if err != nil {
			return err
		}
		fileInfoStr, err := json.MarshalIndent(fileInfo, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(fileInfoStr))
	case "load":
		fileName := gFlags.Get("-a").GetStringValue(pArgs)
		stream, err := stream.BuildStream(
			pCtx,
			uint64(retryNum),
			inputPath,
			friend,
			hlkClient,
			fileName,
			func(_ []byte, p uint64, s uint64) {
				fmt.Printf("download:[ %dB / %dB ]\n", p, s)
			},
		)
		if err != nil {
			return err
		}
		dstFile := filepath.Join(inputPath, fileName)
		if err := copyFile(dstFile, stream); err != nil {
			return err
		}
	default:
		return errors.Join(ErrUnknownAction, errors.New(do)) // nolint:err113
	}

	return nil
}

func copyFile(dst string, src io.Reader) error {
	destinationFile, err := os.Create(dst) // nolint:gosec
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destinationFile.Close() // nolint:errcheck
	_, err = io.Copy(destinationFile, src)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}
	err = destinationFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}
	return nil
}
