package main

import (
	"context"
	"fmt"
	"os"

	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/projects/chat/pkg/app"
	"github.com/number571/hidden-lake/internal/projects/chat/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/internal/utils/help"
)

var (
	gFlags = flag.NewFlagsBuilder(
		flag.NewFlagBuilder("-v", "--version").
			WithDescription("print version of service"),
		flag.NewFlagBuilder("-h", "--help").
			WithDescription("print information about service"),
		flag.NewFlagBuilder("-p", "--path").
			WithDescription("set path to config, database files").
			WithDefinedValue("."),
		flag.NewFlagBuilder("-n", "--network").
			WithDescription("set network key of connections from build").
			WithDefinedValue("oi4r9NW9Le7fKF9d"),
		flag.NewFlagBuilder("-d", "--database").
			WithDescription("set database local file").
			WithDefinedValue(""),
	).Build()
)

func main() {
	args := os.Args[1:]
	if ok := gFlags.Validate(args); !ok {
		panic("args invalid")
	}

	if gFlags.Get("-v").GetBoolValue(args) {
		fmt.Println(build.GetVersion())
		return
	}

	if gFlags.Get("-h").GetBoolValue(args) {
		help.Println(settings.GetAppName(), settings.CProjectDescription, gFlags)
		return
	}

	var (
		nk = gFlags.Get("-n").GetStringValue(args)
		df = gFlags.Get("-d").GetStringValue(args)
	)

	ctx := context.Background()
	if err := app.NewApp(nk, df).Run(ctx); err != nil {
		panic(err)
	}
}
