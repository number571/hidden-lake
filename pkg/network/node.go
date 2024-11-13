package network

import (
	"context"
	"sync"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/anonymity/queue"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/connkeeper"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	hiddenlake "github.com/number571/hidden-lake"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"
)

var (
	_ IHiddenLakeNode = &sHiddenLakeNode{}
)

type sHiddenLakeNode struct {
	fAnonNode   anonymity.INode
	fConnKeeper connkeeper.IConnKeeper
}

func NewHiddenLakeNode(
	pSettings ISettings,
	pPrivKey asymmetric.IPrivKey,
	pKVDatabase database.IKVDatabase,
	pConnsGetter func() []string,
	pHandlerF IHandlerF,
) IHiddenLakeNode {
	return NewRawHiddenLakeNode(
		anonymity.NewNode(
			anonymity.NewSettings(&anonymity.SSettings{
				FNetworkMask:  hiddenlake.GSettings.FProtoMask.FNetwork,
				FServiceName:  pSettings.GetServiceName(),
				FFetchTimeout: pSettings.GetFetchTimeout(),
			}),
			// Insecure to use logging in real anonymity projects!
			// Logging should only be used in overview or testing;
			pSettings.GetLogger(),
			pKVDatabase,
			network.NewNode(
				network.NewSettings(&network.SSettings{
					FAddress:      pSettings.GetTCPAddress(),
					FMaxConnects:  hiddenlake.GSettings.FNetworkManager.FConnectsLimiter,
					FReadTimeout:  hiddenlake.GSettings.GetReadTimeout(),
					FWriteTimeout: hiddenlake.GSettings.GetWriteTimeout(),
					FConnSettings: conn.NewSettings(&conn.SSettings{
						FMessageSettings:       pSettings.GetMessageSettings(),
						FLimitMessageSizeBytes: pSettings.GetMessageSizeBytes(),
						FWaitReadTimeout:       hiddenlake.GSettings.GetWaitTimeout(),
						FDialTimeout:           hiddenlake.GSettings.GetDialTimeout(),
						FReadTimeout:           hiddenlake.GSettings.GetReadTimeout(),
						FWriteTimeout:          hiddenlake.GSettings.GetWriteTimeout(),
					}),
				}),
				cache.NewLRUCache(hiddenlake.GSettings.FNetworkManager.FCacheHashesCap),
			),
			queue.NewQBProblemProcessor(
				queue.NewSettings(&queue.SSettings{
					FMessageConstructSettings: net_message.NewConstructSettings(&net_message.SConstructSettings{
						FSettings: pSettings.GetMessageSettings(),
						FParallel: pSettings.GetParallel(),
					}),
					FNetworkMask: hiddenlake.GSettings.FProtoMask.FNetwork,
					FQueuePeriod: pSettings.GetQueuePeriod(),
					FPoolCapacity: [2]uint64{
						hiddenlake.GSettings.FQueueCapacity.FMain,
						hiddenlake.GSettings.FQueueCapacity.FRand,
					},
				}),
				func() client.IClient {
					client := client.NewClient(pPrivKey, pSettings.GetMessageSizeBytes())
					if client.GetPayloadLimit() <= encoding.CSizeUint64 {
						panic(`client.GetPayloadLimit() <= encoding.CSizeUint64`)
					}
					return client
				}(),
			),
			asymmetric.NewMapPubKeys(),
		),
		pConnsGetter,
		pHandlerF,
	)
}

func NewRawHiddenLakeNode(
	pOriginNode anonymity.INode,
	pConnsGetter func() []string,
	pHandlerF IHandlerF,
) IHiddenLakeNode {
	return &sHiddenLakeNode{
		fAnonNode: pOriginNode.HandleFunc(
			hiddenlake.GSettings.FProtoMask.FService,
			RequestHandler(pHandlerF),
		),
		fConnKeeper: connkeeper.NewConnKeeper(
			connkeeper.NewSettings(&connkeeper.SSettings{
				FDuration:    hiddenlake.GSettings.GetKeeperPeriod(),
				FConnections: pConnsGetter,
			}),
			pOriginNode.GetNetworkNode(),
		),
	}
}

func (p *sHiddenLakeNode) Run(ctx context.Context) error {
	chCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	const N = 2

	errs := [N]error{}
	wg := sync.WaitGroup{}
	wg.Add(N)

	go func() {
		defer func() { wg.Done(); cancel() }()
		errs[0] = p.fConnKeeper.Run(chCtx)
	}()
	go func() {
		defer func() { wg.Done(); cancel() }()
		errs[1] = p.fAnonNode.Run(chCtx)
	}()

	wg.Wait()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		for _, err := range errs {
			if err != nil {
				return err
			}
		}
		panic("closed without errors")
	}
}

func (p *sHiddenLakeNode) SendRequest(
	pCtx context.Context,
	pPubKey asymmetric.IPubKey,
	pRequest request.IRequest,
) error {
	return p.fAnonNode.SendPayload(
		pCtx,
		pPubKey,
		payload.NewPayload64(
			uint64(hiddenlake.GSettings.FProtoMask.FService),
			pRequest.ToBytes(),
		),
	)
}

func (p *sHiddenLakeNode) FetchRequest(
	pCtx context.Context,
	pPubKey asymmetric.IPubKey,
	pRequest request.IRequest,
) (response.IResponse, error) {
	rspBytes, err := p.fAnonNode.FetchPayload(
		pCtx,
		pPubKey,
		payload.NewPayload32(
			hiddenlake.GSettings.FProtoMask.FService,
			pRequest.ToBytes(),
		),
	)
	if err != nil {
		return nil, err
	}
	return response.LoadResponse(rspBytes)
}

func (p *sHiddenLakeNode) GetOrigNode() anonymity.INode {
	return p.fAnonNode
}
