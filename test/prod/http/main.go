package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2/hybrid"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/pkg/network"
	"github.com/number571/hidden-lake/pkg/network/adapters"
	"github.com/number571/hidden-lake/pkg/network/adapters/http"
	"github.com/number571/hidden-lake/pkg/network/request"
	"github.com/number571/hidden-lake/pkg/network/response"

	hla_http_settings "github.com/number571/hidden-lake/internal/adapters/http/pkg/settings"
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

	retries, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(lenNetworks)

	for networkKey := range networks {
		go func(nk string) {
			defer wg.Done()

			networkByKey, _ := build.GetNetwork(networkKey)
			connections := networkByKey.FConnections.GetByScheme(hla_http_settings.CAppAdapterName)
			if len(connections) == 0 {
				return // pass another adapter
			}

			respTime, err := doTestRequest(nk, retries)
			if err != nil {
				log.Printf("%s: network '%s' has error: %s", hla_http_settings.CAppAdapterName, nk, err.Error())
				return
			}
			log.Printf("%s: network '%s' is working successfully; response time %s", hla_http_settings.CAppAdapterName, nk, respTime)
		}(networkKey)
	}

	wg.Wait()
}

func doTestRequest(networkKey string, retries int) (time.Duration, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		node1, key1 = newNode(networkKey, "node1")
		node2, key2 = newNode(networkKey, "node2")
	)

	go func() { _ = node1.Run(ctx) }()
	go func() { _ = node2.Run(ctx) }()

	_, pKey := exchangeKeys(node1, node2, key1, key2)
	startTime := time.Now()

	msg := "hello, world!"
	for i := 0; i < retries; i++ {
		rsp, err := node1.FetchRequest(
			ctx,
			pKey,
			request.NewRequestBuilder().WithBody([]byte(msg)).Build(),
		)
		if err != nil {
			return 0, err
		}
		if string(rsp.GetBody()) != fmt.Sprintf(echoTemplate, msg) {
			return 0, errors.New("got invalid response") // nolint: err113
		}
	}

	return time.Since(startTime), nil
}

func newNode(networkKey string, name string) (network.IHiddenLakeNode, layer2.IParticipantKey) {
	privKey := asymmetric.NewPrivKey()
	adapterSettings := adapters.NewSettingsByNetworkKey(networkKey)
	node, err := network.NewHiddenLakeNode(
		network.NewSettings(&network.SSettings{
			FAdapterSettings: adapterSettings,
		}),
		func() layer2.IScheme {
			scheme, _ := hybrid.NewScheme(privKey, adapterSettings.GetMessageSizeBytes())
			return scheme
		}(),
		layer2.NewKeysContainer(),
		func() database.IKVDatabase {
			kv, err := database.NewKVDatabase(name + "_" + networkKey + ".db")
			if err != nil {
				panic(err)
			}
			return kv
		}(),
		http.NewHTTPAdapter(
			http.NewSettings(&http.SSettings{
				FAdapterSettings: adapterSettings,
				FServeSettings:   &http.SServeSettings{FSubscribeID: name},
			}),
			cache.NewLRUCache(build.GetSettings().FStorageManager.FCacheHashesCap),
			func() []string {
				networkByKey, _ := build.GetNetwork(networkKey)
				return networkByKey.FConnections.GetByScheme(hla_http_settings.CAppAdapterName)
			},
		),
		func(_ context.Context, _ layer2.IParticipantKey, r request.IRequest) (response.IResponse, error) {
			rsp := []byte(fmt.Sprintf(echoTemplate, string(r.GetBody())))
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
