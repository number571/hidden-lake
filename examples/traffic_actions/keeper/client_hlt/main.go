package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	hlt_client "github.com/number571/hidden-lake/internal/helpers/traffic/pkg/client"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

const (
	cPldHead = 0x1
)

const (
	cLocalAddressHLT = "localhost:9582"
)

func main() {
	ctx := context.Background()

	netSett := net_message.NewConstructSettings(&net_message.SConstructSettings{
		FSettings: net_message.NewSettings(&net_message.SSettings{
			FWorkSizeBits: 22,
			FNetworkKey:   "j2BR39JfDf7Bajx3",
		}),
	})

	readPrivKey, err := os.ReadFile("../_keys/priv_node1.key")
	if err != nil {
		panic(err)
	}

	privKey := asymmetric.LoadPrivKey(string(readPrivKey))
	client := client.NewClient(privKey, (8 << 10))

	if len(os.Args) < 2 {
		panic("len os.Args < 2")
	}

	args := os.Args[1:]
	addr := cLocalAddressHLT

	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			"http://"+addr,
			&http.Client{Timeout: time.Minute},
			netSett.GetSettings(),
		),
	)

	switch args[0] {
	case "w", "write":
		if len(args) != 2 {
			panic("len write.args != 2")
		}

		msg, err := client.EncryptMessage(
			privKey.GetPubKey(),
			payload.NewPayload64(cPldHead, []byte(args[1])).ToBytes(),
		)
		if err != nil {
			panic(err)
		}

		netMsg := net_message.NewMessage(
			netSett,
			payload.NewPayload32(hls_settings.CNetworkMask, msg),
		)

		if err := hltClient.PutMessage(ctx, netMsg); err != nil {
			panic(err)
		}

		fmt.Printf("%x\n", netMsg.GetHash())
	case "r", "read":
		if len(args) != 2 {
			panic("len read.args != 2")
		}

		netMsg, err := hltClient.GetMessage(ctx, args[1])
		if err != nil {
			panic(err)
		}

		if netMsg.GetPayload().GetHead() != hls_settings.CNetworkMask {
			panic("net.payload.head is invalid")
		}

		pubKey, decMsg, err := client.DecryptMessage(
			asymmetric.NewMapPubKeys(privKey.GetPubKey()),
			netMsg.GetPayload().GetBody(),
		)
		if err != nil {
			panic(err)
		}

		pld := payload.LoadPayload64(decMsg)
		if pld == nil {
			panic("payload = nil")
		}

		if pld.GetHead() != cPldHead {
			panic("payload.head != set.head")
		}

		if !bytes.Equal(pubKey.ToBytes(), client.GetPrivKey().GetPubKey().ToBytes()) {
			panic("public key is incorrect")
		}

		fmt.Println(string(pld.GetBody()))
	case "h", "hashes":
		for i := uint64(0); ; i++ {
			hash, err := hltClient.GetHash(ctx, i)
			if err != nil {
				return
			}
			fmt.Printf("[%d] %s\n", i+1, hash)
		}
	default:
		panic("unknown mode")
	}
}
