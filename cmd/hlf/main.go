package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/applications/filesharer/pkg/app"
	"github.com/number571/hidden-lake/internal/applications/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/internal/utils/help"
)

var (
	gFlags = flag.NewFlagsBuilder(
		flag.NewFlagBuilder("v", "version").
			WithDescription("print information about service"),
		flag.NewFlagBuilder("h", "help").
			WithDescription("print version of service"),
		flag.NewFlagBuilder("p", "path").
			WithDescription("set path to config, database files").
			WithDefaultValue("."),
	).Build()
)

func main() {
	args := os.Args[1:]

	if ok := gFlags.Validate(args); !ok {
		log.Fatal("args invalid")
		return
	}

	if gFlags.Get("version").GetBoolValue(args) {
		fmt.Println(build.GVersion)
		return
	}

	if gFlags.Get("help").GetBoolValue(args) {
		help.Println(settings.CServiceFullName, settings.CServiceDescription, gFlags)
		return
	}

	app, err := app.InitApp(args, gFlags)
	if err != nil {
		panic(err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	closed := make(chan struct{})
	defer func() {
		cancel()
		<-closed
	}()

	go func() {
		defer func() { closed <- struct{}{} }()
		if err := app.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			log.Fatal(err)
		}
	}()

	<-shutdown
}
