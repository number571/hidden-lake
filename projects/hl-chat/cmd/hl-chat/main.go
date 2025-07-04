package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/number571/hidden-lake/projects/hl-chat/internal/app"
)

func main() {
	pv := flag.Bool("v", false, "print version")
	nk := flag.String("n", "oi4r9NW9Le7fKF9d", "network key")
	df := flag.String("f", "hl-chat.db", "database file")
	flag.Parse()

	if *pv {
		fmt.Println("v0.0.1")
		return
	}

	ctx := context.Background()
	if err := app.NewApp(*nk, *df).Run(ctx); err != nil {
		panic(err)
	}
}
