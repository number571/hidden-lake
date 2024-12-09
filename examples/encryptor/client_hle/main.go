package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/hidden-lake/build"
	hle_client "github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/client"
)

const (
	cLocalAddressHLE = "localhost:8551"
)

func main() {
	ctx := context.Background()

	if len(os.Args) != 3 {
		panic("len(os.Args) != 3")
	}

	netSett := net_message.NewSettings(&net_message.SSettings{
		FWorkSizeBits: 22,
		FNetworkKey:   "j2BR39JfDf7Bajx3",
	})

	hleClient := hle_client.NewClient(
		hle_client.NewBuilder(),
		hle_client.NewRequester(
			cLocalAddressHLE,
			&http.Client{Timeout: time.Minute},
			netSett,
		),
	)

	switch os.Args[1] {
	case "e", "encrypt":
		netMsg, err := hleClient.EncryptMessage(
			ctx,
			"IAM",
			payload.NewPayload64(uint64(build.GSettings.FProtoMask.FService), []byte(os.Args[2])),
		)
		if err != nil {
			panic(err)
		}

		fmt.Println(netMsg.ToString())
	case "d", "decrypt":
		netMsg, err := net_message.LoadMessage(netSett, string(os.Args[2]))
		if err != nil {
			panic(err)
		}

		_, data, err := hleClient.DecryptMessage(ctx, netMsg)
		if err != nil {
			panic(err)
		}

		if data.GetHead() != uint64(build.GSettings.FProtoMask.FService) {
			panic("service mask error")
		}

		fmt.Println(string(data.GetBody()))
	default:
		panic("unknown mode")
	}
}
