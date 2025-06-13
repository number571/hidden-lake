package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/pkg/adapters"
	"github.com/number571/hidden-lake/pkg/adapters/tcp"
	"github.com/number571/hidden-lake/pkg/network"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"

	hla_tcp_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
)

const (
	echoTemplate = "echo: %s;"
)

func main() {
	networks := build.GetNetworks()
	delete(networks, build.CDefaultNetwork)

	lenNetworks := len(networks)
	if lenNetworks == 0 {
		panic("networks is null")
	}

	wg := &sync.WaitGroup{}
	wg.Add(lenNetworks)

	for networkKey := range networks {
		go func(nk string) {
			defer wg.Done()
			respTime, err := doTestRequest(nk)
			if err != nil {
				log.Printf("network '%s' has error: %s", nk, err.Error())
				return
			}
			log.Printf("network '%s' is working successfully; response time %s", nk, respTime)
		}(networkKey)
	}

	wg.Wait()
}

func doTestRequest(networkKey string) (time.Duration, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		node1 = newNode(networkKey, "node1")
		node2 = newNode(networkKey, "node2")
	)

	go func() { _ = node1.Run(ctx) }()
	go func() { _ = node2.Run(ctx) }()

	_, pubKey := exchangeKeys(node1, node2)
	startTime := time.Now()

	msg := "hello, world!"
	rsp, err := node1.FetchRequest(
		ctx,
		pubKey,
		request.NewRequestBuilder().WithBody([]byte(msg)).Build(),
	)
	if err != nil {
		return 0, err
	}

	if string(rsp.GetBody()) != fmt.Sprintf(echoTemplate, msg) {
		return 0, errors.New("got invalid response") // nolint: err113
	}

	return time.Since(startTime), nil
}

func newNode(networkKey string, name string) network.IHiddenLakeNode {
	adapterSettings := adapters.NewSettingsByNetworkKey(networkKey)
	return network.NewHiddenLakeNode(
		network.NewSettings(&network.SSettings{
			FAdapterSettings: adapterSettings,
		}),
		asymmetric.NewPrivKey(),
		func() database.IKVDatabase {
			kv, err := database.NewKVDatabase(name + "_" + networkKey + ".db")
			if err != nil {
				panic(err)
			}
			return kv
		}(),
		tcp.NewTCPAdapter(
			tcp.NewSettings(&tcp.SSettings{
				FAdapterSettings: adapterSettings,
			}),
			cache.NewLRUCache(build.GetSettings().FStorageManager.FCacheHashesCap),
			func() []string {
				networkByKey, _ := build.GetNetwork(networkKey)
				return networkByKey.FConnections.GetByScheme(hla_tcp_settings.CServiceAdapterScheme)
			},
		),
		func(_ context.Context, _ asymmetric.IPubKey, r request.IRequest) (response.IResponse, error) {
			rsp := []byte(fmt.Sprintf(echoTemplate, string(r.GetBody())))
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
