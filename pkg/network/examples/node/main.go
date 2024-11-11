package main

import (
	"context"
	"fmt"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/database"
	hiddenlake "github.com/number571/hidden-lake"
	"github.com/number571/hidden-lake/pkg/network"
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
		node2 = runNode(ctx, "node2", networkKey).HandleFunc(
			serviceMask,
			func(_ context.Context, _ anonymity.INode, _ asymmetric.IPubKey, b []byte) ([]byte, error) {
				return []byte(fmt.Sprintf("echo: %s", string(b))), nil
			},
		)
	)

	_, pubKey := exchangeKeys(node1, node2)
	rsp, _ := node1.FetchPayload(
		ctx,
		pubKey,
		payload.NewPayload32(serviceMask, []byte("hello, world!")),
	)

	fmt.Println(string(rsp))
}

func runNode(ctx context.Context, name, networkKey string) anonymity.INode {
	node := network.NewHiddenLakeNode(
		name,
		network.NewSettingsByNetworkKey(networkKey, nil),
		asymmetric.NewPrivKey(),
		func() database.IKVDatabase {
			kv, _ := database.NewKVDatabase(name + ".db")
			return kv
		}(),
	)
	network := hiddenlake.GNetworks[networkKey]
	for _, c := range network.FConnections {
		_ = node.GetNetworkNode().AddConnection(
			ctx,
			fmt.Sprintf("%s:%d", c.FHost, c.FPort),
		)
	}
	go func() { _ = node.Run(ctx) }()
	return node
}

func exchangeKeys(node1, node2 anonymity.INode) (asymmetric.IPubKey, asymmetric.IPubKey) {
	pubKey1 := node1.GetMessageQueue().GetClient().GetPrivKey().GetPubKey()
	pubKey2 := node2.GetMessageQueue().GetClient().GetPrivKey().GetPubKey()

	node1.GetMapPubKeys().SetPubKey(pubKey2)
	node2.GetMapPubKeys().SetPubKey(pubKey1)

	return pubKey1, pubKey2
}
