package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/pkg/adapters"
	"github.com/number571/hidden-lake/pkg/adapters/tcp"
	"github.com/number571/hidden-lake/pkg/network"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"

	anon_logger "github.com/number571/go-peer/pkg/anonymity/logger"
)

const (
	relayerAddress = "localhost:9999"
	msgSizeBytes   = uint64(8 << 10)
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		node1   = newNode(ctx, "node1")
		node2   = newNode(ctx, "node2")
		adapter = newTCPAdapter(relayerAddress, nil)
	)

	go func() { _ = adapter.Run(ctx) }()
	go func() { _ = runAsRelayer(ctx, adapter) }()

	go func() { _ = node1.Run(ctx) }()
	go func() { _ = node2.Run(ctx) }()

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
			FQueuePeriod:  time.Second,
			FFetchTimeout: time.Minute,
			FAdapterSettings: adapters.NewSettings(&adapters.SSettings{
				FMessageSizeBytes: msgSizeBytes,
			}),
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
		newTCPAdapter("", []string{relayerAddress}),
		func(_ context.Context, _ asymmetric.IPubKey, r request.IRequest) (response.IResponse, error) {
			rsp := []byte(fmt.Sprintf("echo: %s", string(r.GetBody())))
			return response.NewResponseBuilder().WithBody(rsp).Build(), nil
		},
	)
}

func newTCPAdapter(addr string, conns []string) adapters.IRunnerAdapter {
	return tcp.NewTCPAdapter(
		tcp.NewSettings(&tcp.SSettings{
			FAddress: addr,
			FAdapterSettings: adapters.NewSettings(&adapters.SSettings{
				FMessageSizeBytes: msgSizeBytes,
			}),
		}),
		cache.NewLRUCache(1<<10),
		func() []string { return conns },
	)
}

func runAsRelayer(ctx context.Context, adapter adapters.IRunnerAdapter) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := adapter.Consume(ctx)
			if err != nil {
				continue
			}
			_ = adapter.Produce(ctx, msg)
		}
	}
}

func exchangeKeys(hlNode1, hlNode2 network.IHiddenLakeNode) (asymmetric.IPubKey, asymmetric.IPubKey) {
	node1 := hlNode1.GetOriginNode()
	node2 := hlNode2.GetOriginNode()

	pubKey1 := node1.GetQBProcessor().GetClient().GetPrivKey().GetPubKey()
	pubKey2 := node2.GetQBProcessor().GetClient().GetPrivKey().GetPubKey()

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
