package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	hiddenlake "github.com/number571/hidden-lake"
	"github.com/number571/hidden-lake/internal/helpers/loader/pkg/app"
	"github.com/number571/hidden-lake/internal/utils/flag"
)

func main() {
	args := os.Args[1:]

	if flag.GetBoolFlagValue(args, []string{"v", "version"}) {
		fmt.Println(hiddenlake.GVersion)
		return
	}

	if flag.GetBoolFlagValue(args, []string{"h", "help"}) {
		fmt.Print(
			"Hidden Lake Loader (HLL)\n" +
				"Description: distributes the stored traffic between nodes\n" +
				"Arguments:\n" +
				"[ -h, --help    ] - print information about service\n" +
				"[ -v, --version ] - print version of service\n" +
				"[ -p, --path    ] - set path to config, database files\n" +
				"[ -n, --network ] - set network key for connections\n",
		)
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
