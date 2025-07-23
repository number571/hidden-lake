package network

import (
	"context"
	"crypto/ed25519"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/cmd/hlp/hlp-chat/internal/request"
	"github.com/number571/hidden-lake/pkg/adapters"
	"github.com/number571/hidden-lake/pkg/adapters/tcp"
	"github.com/number571/hidden-lake/pkg/network"
	"github.com/number571/hidden-lake/pkg/response"

	pkg_request "github.com/number571/hidden-lake/pkg/request"
)

var (
	_ IHiddenLakeChatNode = &sHiddenLakeChatNode{}
)

type sHiddenLakeChatNode struct {
	fChannelKey     asymmetric.IPubKey
	fPrivateKey     ed25519.PrivateKey
	fHiddenLakeNode network.IHiddenLakeNode
}

func NewHiddenLakeChatNode(
	pNetworkKey string,
	pDatabase database.IKVDatabase,
	pChanPrivKey asymmetric.IPrivKey,
	pPrivKey ed25519.PrivateKey,
	pCallbackFunc ICallbackFunc,
) IHiddenLakeChatNode {
	networkByKey, ok := build.GetNetwork(pNetworkKey)
	if !ok {
		panic("network key undefined")
	}

	adapterSettings := adapters.NewSettingsByNetworkKey(pNetworkKey)
	connections := networkByKey.FConnections.GetByScheme("tcp")

	chanKey := pChanPrivKey.GetPubKey()
	node := network.NewHiddenLakeNode(
		network.NewSettings(&network.SSettings{
			FAdapterSettings: adapterSettings,
		}),
		pChanPrivKey,
		pDatabase,
		tcp.NewTCPAdapter(
			tcp.NewSettings(&tcp.SSettings{
				FAdapterSettings: adapterSettings,
			}),
			cache.NewLRUCache(2<<10),
			func() []string { return connections },
		),
		func(_ context.Context, pk asymmetric.IPubKey, r pkg_request.IRequest) (response.IResponse, error) {
			pubKey, hash, ok := request.ValidateRequest(chanKey, r)
			if ok {
				pCallbackFunc(pubKey, hash, string(r.GetBody()))
			}
			return nil, nil
		},
	)

	node.GetOriginNode().GetMapPubKeys().SetPubKey(chanKey)
	return &sHiddenLakeChatNode{
		fChannelKey:     chanKey,
		fPrivateKey:     pPrivKey,
		fHiddenLakeNode: node,
	}
}

func (p *sHiddenLakeChatNode) Run(pCtx context.Context) error {
	return p.fHiddenLakeNode.GetOriginNode().Run(pCtx)
}

func (p *sHiddenLakeChatNode) GetMessageLimitSize() uint64 {
	pldLimit := p.fHiddenLakeNode.GetOriginNode().GetQBProcessor().GetClient().GetPayloadLimit()
	return request.GetMessageLimitSize(pldLimit)
}

func (p *sHiddenLakeChatNode) SendMessage(pCtx context.Context, pMsg string) error {
	return p.fHiddenLakeNode.SendRequest(
		pCtx,
		p.fChannelKey,
		request.BuildRequest(p.fChannelKey, p.fPrivateKey, pMsg),
	)
}
