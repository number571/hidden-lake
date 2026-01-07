package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/number571/hidden-lake/build"
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	hls_pinger_client "github.com/number571/hidden-lake/internal/services/pinger/pkg/client"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/internal/utils/help"
)

var (
	gFlags = flag.NewFlagsBuilder(
		flag.NewFlagBuilder("-v", "--version").
			WithDescription("print version of application"),
		flag.NewFlagBuilder("-h", "--help").
			WithDescription("print information about application"),
		flag.NewFlagBuilder("-k", "--kernel").
			WithDescription("set internal address of the HLK").
			WithDefinedValue("localhost:9572"),
		flag.NewFlagBuilder("-f", "--friend").
			WithDescription("set alias name of the friend").
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

	fmt.Println("ok")
}

func runFunction(pCtx context.Context, pArgs []string) error {
	hlkClient := hlk_client.NewClient(
		hlk_client.NewBuilder(),
		hlk_client.NewRequester(
			gFlags.Get("-k").GetStringValue(pArgs),
			&http.Client{Timeout: build.GetSettings().GetHttpCallbackTimeout()},
		),
	)
	return hls_pinger_client.NewClient(
		hls_pinger_client.NewBuilder(),
		hls_pinger_client.NewRequester(hlkClient),
	).Ping(pCtx, gFlags.Get("-f").GetStringValue(pArgs))
}
