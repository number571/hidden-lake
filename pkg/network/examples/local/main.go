package main

import (
	"context"
	"fmt"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/pkg/network"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"
)

const (
	nodeTCPAddress = "localhost:9999"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		node1 = runNode(ctx, "node1", nodeTCPAddress, func() []string { return nil })
		node2 = runNode(ctx, "node2", "", func() []string { return []string{nodeTCPAddress} })
	)

	_, pubKey := exchangeKeys(node1, node2)
	rsp, _ := node1.FetchRequest(
		ctx,
		pubKey,
		request.NewRequest().WithBody([]byte("hello, world!")),
	)

	fmt.Println(string(rsp.GetBody()))
}

func runNode(
	ctx context.Context,
	dbPath, tcpAddr string,
	connsF func() []string,
) network.IHiddenLakeNode {
	const networkKey = "custom_network_key"
	node := network.NewHiddenLakeNode(
		network.NewSettings(&network.SSettings{
			FMessageSettings: message.NewSettings(&message.SSettings{
				FWorkSizeBits: 10,
				FNetworkKey:   networkKey,
			}),
			FQueuePeriod:      time.Second,
			FFetchTimeout:     time.Minute,
			FMessageSizeBytes: (8 << 10),
			SSubSettings: network.SSubSettings{
				FTCPAddress: tcpAddr,
			},
		}),
		asymmetric.NewPrivKey(),
		func() database.IKVDatabase {
			kv, _ := database.NewKVDatabase(dbPath + ".db")
			return kv
		}(),
		connsF,
		func(_ context.Context, _ asymmetric.IPubKey, r request.IRequest) (response.IResponse, error) {
			rsp := []byte(fmt.Sprintf("echo: %s", string(r.GetBody())))
			return response.NewResponse().WithBody(rsp), nil
		},
	)
	go func() { _ = node.Run(ctx) }()
	go func() { _ = node.GetOrigNode().GetNetworkNode().Listen(ctx) }()
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
