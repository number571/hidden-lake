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
	"github.com/number571/hidden-lake/internal/helpers/traffic/pkg/app"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/internal/utils/help"

	_ "embed"
)

var (
	//go:embed help.yml
	gHelpYaml []byte
)

func main() {
	args := os.Args[1:]

	if flag.GetBoolFlagValue(args, []string{"v", "version"}) {
		fmt.Println(build.GVersion)
		return
	}

	if flag.GetBoolFlagValue(args, []string{"h", "help"}) {
		help.Println(gHelpYaml)
		return
	}

	app, err := app.InitApp(args)
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
