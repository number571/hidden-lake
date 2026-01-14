package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/number571/hidden-lake/build"
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
		flag.NewFlagBuilder("-s", "--service").
			WithDescription("set internal address of the HLS").
			WithDefinedValue("localhost:9551"),
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
	return hls_pinger_client.NewClient(
		hls_pinger_client.NewRequester(
			gFlags.Get("-s").GetStringValue(pArgs),
			&http.Client{Timeout: build.GetSettings().GetHttpCallbackTimeout()},
		),
	).PingFriend(pCtx, gFlags.Get("-f").GetStringValue(pArgs))
}
