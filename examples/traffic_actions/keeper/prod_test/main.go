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
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/encoding"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	hiddenlake "github.com/number571/hidden-lake"
	hlt_client "github.com/number571/hidden-lake/internal/helpers/traffic/pkg/client"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

const (
	cPldHead     = 0x1
	cKeySize     = 4096
	cMsgSize     = (8 << 10)
	cWrkSize     = 22
	cPrivKeyPath = "../_keys/priv_node1.key"
)

func main() {
	ctx := context.Background()
	randString := random.NewRandom().GetString(16)

	readPrivKey, err := os.ReadFile(cPrivKeyPath)
	if err != nil {
		panic(err)
	}

	privKey := asymmetric.LoadPrivKey(string(readPrivKey))
	client := client.NewClient(privKey, cMsgSize)

	i := 0
	for key, network := range hiddenlake.GNetworks {
		connect := network.FConnections[0].FHost + ":9582" // HTTP port

		netSett := net_message.NewConstructSettings(&net_message.SConstructSettings{
			FSettings: net_message.NewSettings(&net_message.SSettings{
				FWorkSizeBits: cWrkSize,
				FNetworkKey:   key,
			}),
		})

		hltClient := hlt_client.NewClient(
			hlt_client.NewBuilder(),
			hlt_client.NewRequester(
				"http://"+connect,
				&http.Client{Timeout: 5 * time.Second},
				netSett.GetSettings(),
			),
		)

		msg, err := client.EncryptMessage(
			privKey.GetPubKey(),
			payload.NewPayload64(cPldHead, []byte(randString)).ToBytes(),
		)
		if err != nil {
			fmt.Printf("%d. %s: %s\n", i+1, connect, err)
			continue
		}

		netMsg := net_message.NewMessage(
			netSett,
			payload.NewPayload32(hls_settings.CNetworkMask, msg),
		)

		start1 := time.Now()
		if err := hltClient.PutMessage(ctx, netMsg); err != nil {
			fmt.Printf("%d. %s: %s\n", i+1, connect, err)
			continue
		}

		start2 := time.Now()
		gotNetMsg, err := hltClient.GetMessage(ctx, encoding.HexEncode(netMsg.GetHash()))
		if err != nil {
			fmt.Printf("%d. %s: %s\n", i+1, connect, err)
			continue
		}

		if !bytes.Equal(netMsg.ToBytes(), gotNetMsg.ToBytes()) {
			fmt.Printf("%d. %s: !bytes.Equal(netMsg.ToBytes(), gotNetMsg.ToBytes())\n", i+1, connect)
			continue
		}

		fmt.Printf(
			"%d. HLT server '%s' is working properly (response_time=%s,%s);\n",
			i+1,
			connect,
			time.Since(start1),
			time.Since(start2),
		)
	}
}
