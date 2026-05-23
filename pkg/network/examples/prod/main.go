package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2/hybrid"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/pkg/network"
	"github.com/number571/hidden-lake/pkg/network/adapters"
	"github.com/number571/hidden-lake/pkg/network/adapters/tcp"
	"github.com/number571/hidden-lake/pkg/network/request"
	"github.com/number571/hidden-lake/pkg/network/response"

	anon_logger "github.com/number571/go-peer/pkg/anonymity/qb/logger"
	hla_tcp_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
)

const (
	networkKey = "oi4r9NW9Le7fKF9d"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		node1, key1 = newNode("node1")
		node2, key2 = newNode("node2")
	)

	go func() { _ = node1.Run(ctx) }()
	go func() { _ = node2.Run(ctx) }()

	_, pKey := exchangeKeys(node1, node2, key1, key2)

	for {
		timeNow := time.Now()
		rsp, err := node1.FetchRequest(
			ctx,
			pKey,
			request.NewRequestBuilder().WithBody([]byte("hello, world!")).Build(),
		)
		if err != nil {
			fmt.Printf("error:(%s)\n", err.Error())
			continue
		}
		fmt.Printf("response[%s]:(%s)\n", time.Since(timeNow), string(rsp.GetBody()))
	}
}

func newNode(name string) (network.IHiddenLakeNode, layer2.IParticipantKey) {
	privKey := asymmetric.NewPrivKey()
	adapterSettings := adapters.NewSettingsByNetworkKey(networkKey)
	node, err := network.NewHiddenLakeNode(
		network.NewSettings(&network.SSettings{
			FAdapterSettings: adapterSettings,
			FServeSettings: &network.SServeSettings{
				FServiceName: name,
				FLogger:      getLogger(),
			},
		}),
		hybrid.NewScheme(privKey, adapterSettings.GetMessageSizeBytes()),
		layer2.NewKeysContainer(),
		func() database.IKVDatabase {
			kv, err := database.NewKVDatabase(name + ".db")
			if err != nil {
				panic(err)
			}
			return kv
		}(),
		tcp.NewTCPAdapter(
			tcp.NewSettings(&tcp.SSettings{
				FAdapterSettings: adapterSettings,
			}),
			cache.NewLRUCache(1<<10),
			func() []string {
				networkByKey, _ := build.GetNetwork(networkKey)
				return networkByKey.FConnections.GetByScheme(hla_tcp_settings.CAppAdapterName)
			},
		),
		func(_ context.Context, _ layer2.IParticipantKey, r request.IRequest) (response.IResponse, error) {
			rsp := []byte("echo: " + string(r.GetBody()))
			return response.NewResponseBuilder().WithBody(rsp).Build(), nil
		},
	)
	if err != nil {
		panic(err)
	}
	return node, privKey.GetPubKey()
}

func exchangeKeys(hlNode1, hlNode2 network.IHiddenLakeNode, key1, key2 layer2.IParticipantKey) (layer2.IParticipantKey, layer2.IParticipantKey) {
	node1 := hlNode1.GetOriginNode()
	node2 := hlNode2.GetOriginNode()

	node1.GetKeysContainer().Add(key2)
	node2.GetKeysContainer().Add(key1)

	return key1, key2
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
