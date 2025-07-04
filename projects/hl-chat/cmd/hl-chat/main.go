package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/number571/hidden-lake/projects/hl-chat/internal/app"
)

func main() {
	pv := flag.Bool("v", false, "print version")
	nk := flag.String("n", "oi4r9NW9Le7fKF9d", "set network key")
	flag.Parse()

	if *pv {
		fmt.Println("v1.0.0")
		return
	}

	ctx := context.Background()
	if err := app.NewApp(*nk).Run(ctx); err != nil {
		panic(err)
	}
}
