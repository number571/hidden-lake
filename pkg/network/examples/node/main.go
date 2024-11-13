package main

import (
	"context"
	"fmt"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/storage/database"
	hiddenlake "github.com/number571/hidden-lake"
	"github.com/number571/hidden-lake/pkg/network"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"
)

const (
	serviceMask = uint32(0x011111111)
	networkKey  = "oi4r9NW9Le7fKF9d"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		node1 = runNode(ctx, "node1", networkKey)
		node2 = runNode(ctx, "node2", networkKey)
	)

	_, pubKey := exchangeKeys(node1, node2)
	rsp, _ := node1.FetchRequest(
		ctx,
		pubKey,
		request.NewRequest().WithBody([]byte("hello, world!")),
	)

	fmt.Println(string(rsp.GetBody()))
}

func runNode(ctx context.Context, dbPath, networkKey string) network.IHiddenLakeNode {
	node := network.NewHiddenLakeNode(
		network.NewSettingsByNetworkKey(networkKey, nil),
		asymmetric.NewPrivKey(),
		func() database.IKVDatabase {
			kv, _ := database.NewKVDatabase(dbPath + ".db")
			return kv
		}(),
		func() []string {
			network := hiddenlake.GNetworks[networkKey]
			conns := make([]string, 0, len(network.FConnections))
			for _, c := range network.FConnections {
				conns = append(conns, fmt.Sprintf("%s:%d", c.FHost, c.FPort))
			}
			return conns
		},
		func(_ context.Context, _ asymmetric.IPubKey, r request.IRequest) (response.IResponse, error) {
			rsp := []byte(fmt.Sprintf("echo: %s", string(r.GetBody())))
			return response.NewResponse().WithBody(rsp), nil
		},
	)
	go func() { _ = node.Run(ctx) }()
	return node
}

func exchangeKeys(hlNode1, hlNode2 network.IHiddenLakeNode) (asymmetric.IPubKey, asymmetric.IPubKey) {
	node1 := hlNode1.GetOrigNode()
	node2 := hlNode2.GetOrigNode()

	pubKey1 := node1.GetMessageQueue().GetClient().GetPrivKey().GetPubKey()
	pubKey2 := node2.GetMessageQueue().GetClient().GetPrivKey().GetPubKey()

	node1.GetMapPubKeys().SetPubKey(pubKey2)
	node2.GetMapPubKeys().SetPubKey(pubKey1)

	return pubKey1, pubKey2
}
