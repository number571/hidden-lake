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
	hls_filesharer_client "github.com/number571/hidden-lake/internal/services/filesharer/pkg/client"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
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
			WithDescription("set path to load files").
			WithDefinedValue("."),
		flag.NewFlagBuilder("-s", "--service").
			WithDescription("set internal address of the HLS").
			WithDefinedValue("localhost:9541"),
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
		help.Println(settings.CAppFullName+"-cli", settings.CAppDescription, gFlags)
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

	hlfClient := hls_filesharer_client.NewClient(
		hls_filesharer_client.NewRequester(
			gFlags.Get("-s").GetStringValue(pArgs),
			&http.Client{Timeout: build.GetSettings().GetHttpCallbackTimeout()},
		),
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
		fmt.Println(serializeJSON(fileInfoList))
	case "info":
		fileName := gFlags.Get("-a").GetStringValue(pArgs)
		fileInfo, err := hlfClient.GetFileInfo(pCtx, friend, fileName)
		if err != nil {
			return err
		}
		fmt.Println(serializeJSON(fileInfo))
	case "load":
		fileName := gFlags.Get("-a").GetStringValue(pArgs)
		dstFile, err := os.OpenFile( // nolint: gosec
			filepath.Join(inputPath, fileName),
			os.O_CREATE|os.O_WRONLY,
			0600,
		)
		if err != nil {
			return err
		}
		pw := &processWriter{fW: dstFile}
		if err := hlfClient.DownloadFile(pw, pCtx, friend, fileName); err != nil {
			return err
		}
		fmt.Printf("\ndone!\n")
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

func serializeJSON(pData interface{}) string {
	res, _ := json.MarshalIndent(pData, "", "\t")
	return string(res)
}
