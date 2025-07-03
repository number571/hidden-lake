package main

import (
	"context"
	"flag"

	"github.com/number571/hidden-lake/projects/hl-chat/internal/app"
)

func main() {
	nk := flag.String("n", "oi4r9NW9Le7fKF9d", "network key of the hidden-lake")
	flag.Parse()

	ctx := context.Background()
	if err := app.NewApp(*nk).Run(ctx); err != nil {
		panic(err)
	}
}
