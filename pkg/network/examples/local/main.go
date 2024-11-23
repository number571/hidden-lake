package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/pkg/network"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"

	gopeer_network "github.com/number571/go-peer/pkg/network"
	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
)

const (
	relayerTCPAddress = "localhost:9999"
	msgSizeBytes      = uint64(8 << 10)
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		node1   = newNode(ctx, "node1")
		node2   = newNode(ctx, "node2")
		relayer = newRelayer()
	)

	go func() { _ = node1.Run(ctx) }()
	go func() { _ = node2.Run(ctx) }()
	go func() { _ = relayer.Listen(ctx) }()

	_, pubKey := exchangeKeys(node1, node2)

	for {
		rsp, err := node1.FetchRequest(
			ctx,
			pubKey,
			request.NewRequestBuilder().WithBody([]byte("hello, world!")).Build(),
		)
		if err != nil {
			fmt.Printf("error:(%s)\n", err.Error())
			continue
		}
		fmt.Printf("response:(%s)\n", string(rsp.GetBody()))
	}
}

func newNode(ctx context.Context, name string) network.IHiddenLakeNode {
	return network.NewHiddenLakeNode(
		network.NewSettings(&network.SSettings{
			FQueuePeriod:      time.Second,
			FFetchTimeout:     time.Minute,
			FMessageSizeBytes: msgSizeBytes,
			FSubSettings: &network.SSubSettings{
				FLogger:      getLogger(),
				FServiceName: name,
			},
		}),
		asymmetric.NewPrivKey(),
		func() database.IKVDatabase {
			kv, _ := database.NewKVDatabase(name + ".db")
			return kv
		}(),
		func() []string { return []string{relayerTCPAddress} },
		func(_ context.Context, _ asymmetric.IPubKey, r request.IRequest) (response.IResponse, error) {
			rsp := []byte(fmt.Sprintf("echo: %s", string(r.GetBody())))
			return response.NewResponseBuilder().WithBody(rsp).Build(), nil
		},
	)
}

func newRelayer() gopeer_network.INode {
	timeout := 5 * time.Second
	return gopeer_network.NewNode(
		gopeer_network.NewSettings(&gopeer_network.SSettings{
			FAddress:      relayerTCPAddress,
			FMaxConnects:  256,
			FReadTimeout:  timeout,
			FWriteTimeout: timeout,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FMessageSettings:       message.NewSettings(&message.SSettings{}),
				FLimitMessageSizeBytes: msgSizeBytes,
				FWaitReadTimeout:       time.Hour,
				FDialTimeout:           timeout,
				FReadTimeout:           timeout,
				FWriteTimeout:          timeout,
			}),
		}),
		cache.NewLRUCache(2048),
	).HandleFunc(
		build.GSettings.FProtoMask.FNetwork,
		func(ctx context.Context, node gopeer_network.INode, _ conn.IConn, msg message.IMessage) error {
			node.BroadcastMessage(ctx, msg)
			return nil
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
			logGetter, ok := ia.(anon_logger.ILogGetter)
			if !ok {
				panic("got invalid log arg")
			}
			return fmt.Sprintf(
				"name=%s code=%02x hash=%X proof=%08d bytes=%d",
				logGetter.GetService(),
				logGetter.GetType(),
				logGetter.GetHash()[:16],
				logGetter.GetProof(),
				logGetter.GetSize(),
			)
		},
	)
}
