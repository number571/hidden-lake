package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/pkg/adapters"
	"github.com/number571/hidden-lake/pkg/adapters/tcp"
	"github.com/number571/hidden-lake/pkg/network"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"

	anon_logger "github.com/number571/go-peer/pkg/anonymity/logger"
)

const (
	networkKey = "oi4r9NW9Le7fKF9d"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		node1 = newNode(ctx, "node1")
		node2 = newNode(ctx, "node2")
	)

	go func() { _ = node1.Run(ctx) }()
	go func() { _ = node2.Run(ctx) }()

	_, pubKey := exchangeKeys(node1, node2)

	for {
		timeNow := time.Now()
		rsp, err := node1.FetchRequest(
			ctx,
			pubKey,
			request.NewRequestBuilder().WithBody([]byte("hello, world!")).Build(),
		)
		if err != nil {
			fmt.Printf("error:(%s)\n", err.Error())
			continue
		}
		fmt.Printf("response[%s]:(%s)\n", time.Since(timeNow), string(rsp.GetBody()))
	}
}

func newNode(_ context.Context, name string) network.IHiddenLakeNode {
	networkByKey, _ := build.GetNetwork(networkKey)
	connects := make([]string, 0, len(networkByKey.FConnections))
	for _, c := range networkByKey.FConnections {
		u, err := url.Parse(c)
		if err != nil || u.Scheme != "tcp" {
			continue
		}
		connects = append(connects, u.Host)
	}
	return network.NewHiddenLakeNode(
		network.NewSettings(&network.SSettings{
			FAdapterSettings: adapters.NewSettingsByNetworkKey(networkKey),
			FQBPSettings: &network.SQBPSettings{
				FFetchTimeout: time.Minute,
				FQueuePeriod:  5 * time.Second,
			},
			FSrvSettings: &network.SSrvSettings{
				FServiceName: name,
				FLogger:      getLogger(),
			},
		}),
		asymmetric.NewPrivKey(),
		func() database.IKVDatabase {
			kv, err := database.NewKVDatabase(name + ".db")
			if err != nil {
				panic(err)
			}
			return kv
		}(),
		tcp.NewTCPAdapter(
			tcp.NewSettings(&tcp.SSettings{
				FAdapterSettings: adapters.NewSettings(&adapters.SSettings{
					FNetworkKey:       networkKey,
					FMessageSizeBytes: networkByKey.FMessageSizeBytes,
					FWorkSizeBits:     networkByKey.FWorkSizeBits,
				}),
			}),
			cache.NewLRUCache(1<<10),
			func() []string { return connects },
		),
		func(_ context.Context, _ asymmetric.IPubKey, r request.IRequest) (response.IResponse, error) {
			rsp := []byte("echo: " + string(r.GetBody()))
			return response.NewResponseBuilder().WithBody(rsp).Build(), nil
		},
	)
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
