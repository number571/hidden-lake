package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/pkg/network"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"

	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
)

const (
	nodeTCPAddress = "localhost:9999"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		node1 = newNode(ctx, "node1", nodeTCPAddress, func() []string { return nil })
		node2 = newNode(ctx, "node2", "", func() []string { return []string{nodeTCPAddress} })
	)

	go func() { _ = node1.Run(ctx) }()
	go func() { _ = node2.Run(ctx) }()
	go func() { _ = node1.GetOriginNode().GetNetworkNode().Listen(ctx) }()

	_, pubKey := exchangeKeys(node1, node2)

	for {
		rsp, err := node1.FetchRequest(
			ctx,
			pubKey,
			request.NewRequest().WithBody([]byte("hello, world!")),
		)
		if err != nil {
			fmt.Printf("error:(%s)\n", err.Error())
			continue
		}
		fmt.Printf("response:(%s)\n", string(rsp.GetBody()))
	}
}

func newNode(
	ctx context.Context,
	name, tcpAddr string,
	connsF func() []string,
) network.IHiddenLakeNode {
	return network.NewHiddenLakeNode(
		network.NewSettings(&network.SSettings{
			FMessageSettings: message.NewSettings(&message.SSettings{
				FWorkSizeBits: 10,
				FNetworkKey:   "custom_network_key",
			}),
			FQueuePeriod:      time.Second,
			FFetchTimeout:     time.Minute,
			FMessageSizeBytes: (8 << 10),
			SSubSettings: network.SSubSettings{
				FTCPAddress:  tcpAddr,
				FLogger:      getLogger(),
				FServiceName: name,
			},
		}),
		asymmetric.NewPrivKey(),
		func() database.IKVDatabase {
			kv, _ := database.NewKVDatabase(name + ".db")
			return kv
		}(),
		connsF,
		func(_ context.Context, _ asymmetric.IPubKey, r request.IRequest) (response.IResponse, error) {
			rsp := []byte(fmt.Sprintf("echo: %s", string(r.GetBody())))
			return response.NewResponse().WithBody(rsp), nil
		},
	)
}

func exchangeKeys(hlNode1, hlNode2 network.IHiddenLakeNode) (asymmetric.IPubKey, asymmetric.IPubKey) {
	node1 := hlNode1.GetOriginNode()
	node2 := hlNode2.GetOriginNode()

	pubKey1 := node1.GetMessageQueue().GetClient().GetPrivKey().GetPubKey()
	pubKey2 := node2.GetMessageQueue().GetClient().GetPrivKey().GetPubKey()

	node1.GetMapPubKeys().SetPubKey(pubKey2)
	node2.GetMapPubKeys().SetPubKey(pubKey1)

	return pubKey1, pubKey2
}

func getLogger() logger.ILogger {
	return logger.NewLogger(
		logger.NewSettings(&logger.SSettings{
			FInfo: os.Stdout,
			FWarn: os.Stdout,
			FErro: os.Stderr,
		}),
		func(ia logger.ILogArg) string {
			logGetterFactory, ok := ia.(anon_logger.ILogGetterFactory)
			if !ok {
				panic("got invalid log arg")
			}
			logGetter := logGetterFactory.Get()
			hash := make([]byte, hashing.CHasherSize)
			if x := logGetter.GetHash(); x != nil {
				copy(hash, x)
			}
			return fmt.Sprintf(
				"name=%s code=%02x hash=%X proof=%08d bytes=%d",
				logGetter.GetService(),
				logGetter.GetType(),
				hash[:16],
				logGetter.GetProof(),
				logGetter.GetSize(),
			)
		},
	)
}
