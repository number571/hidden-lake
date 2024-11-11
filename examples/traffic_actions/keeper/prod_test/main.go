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
)

const (
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

	i := 0
	for key, network := range hiddenlake.GNetworks {
		i++
		if key == hiddenlake.CDefaultNetwork {
			continue
		}

		connect := fmt.Sprintf("%s:%d", network.FConnections[0].FHost, 9582)
		netSett := net_message.NewConstructSettings(&net_message.SConstructSettings{
			FSettings: net_message.NewSettings(&net_message.SSettings{
				FWorkSizeBits: network.FWorkSizeBits,
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

		client := client.NewClient(privKey, network.FMessageSizeBytes)
		msg, err := client.EncryptMessage(
			privKey.GetPubKey(),
			payload.NewPayload64(0x01, []byte(randString)).ToBytes(),
		)
		if err != nil {
			fmt.Printf("%d. %s: %s\n", i, connect, err)
			continue
		}

		netMsg := net_message.NewMessage(
			netSett,
			payload.NewPayload32(hiddenlake.GSettings.FProtoMask.FService, msg),
		)

		start1 := time.Now()
		if err := hltClient.PutMessage(ctx, netMsg); err != nil {
			fmt.Printf("%d. %s: %s\n", i, connect, err)
			continue
		}

		start2 := time.Now()
		gotNetMsg, err := hltClient.GetMessage(ctx, encoding.HexEncode(netMsg.GetHash()))
		if err != nil {
			fmt.Printf("%d. %s: %s\n", i, connect, err)
			continue
		}

		if !bytes.Equal(netMsg.ToBytes(), gotNetMsg.ToBytes()) {
			fmt.Printf("%d. %s: !bytes.Equal(netMsg.ToBytes(), gotNetMsg.ToBytes())\n", i, connect)
			continue
		}

		fmt.Printf(
			"%d. HLT server '%s' is working properly (response_time=%s,%s);\n",
			i,
			connect,
			time.Since(start1),
			time.Since(start2),
		)
	}
}
