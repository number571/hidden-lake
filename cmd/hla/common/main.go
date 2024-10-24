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
	"github.com/number571/hidden-lake/internal/adapters/common/pkg/app"
	"github.com/number571/hidden-lake/internal/utils/flag"
)

func main() {
	args := os.Args[1:]

	if flag.GetBoolFlagValue(args, "version") {
		fmt.Println(hiddenlake.CVersion)
		return
	}

	app, err := app.InitApp(args, ".")
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
