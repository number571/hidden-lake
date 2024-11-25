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
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/pkg/handler"
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
	pHandlerF handler.IHandlerF,
) IHiddenLakeNode {
	msgSettings := net_message.NewSettings(&net_message.SSettings{
		FWorkSizeBits: pSettings.GetWorkSizeBits(),
		FNetworkKey:   pSettings.GetNetworkKey(),
	})
	node := anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
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
				FMaxConnects:  build.GSettings.FNetworkManager.FConnectsLimiter,
				FReadTimeout:  build.GSettings.GetReadTimeout(),
				FWriteTimeout: build.GSettings.GetWriteTimeout(),
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FMessageSettings:       msgSettings,
					FLimitMessageSizeBytes: pSettings.GetMessageSizeBytes(),
					FWaitReadTimeout:       build.GSettings.GetWaitTimeout(),
					FDialTimeout:           build.GSettings.GetDialTimeout(),
					FReadTimeout:           build.GSettings.GetReadTimeout(),
					FWriteTimeout:          build.GSettings.GetWriteTimeout(),
				}),
			}),
			cache.NewLRUCache(build.GSettings.FNetworkManager.FCacheHashesCap),
		),
		queue.NewQBProblemProcessor(
			queue.NewSettings(&queue.SSettings{
				FMessageConstructSettings: net_message.NewConstructSettings(&net_message.SConstructSettings{
					FSettings: msgSettings,
					FParallel: pSettings.GetParallel(),
				}),
				FQueuePeriod:  pSettings.GetQueuePeriod(),
				FNetworkMask:  build.GSettings.FProtoMask.FNetwork,
				FConsumersCap: build.GSettings.FQueueProblem.FConsumersCap,
				FQueuePoolCap: [2]uint64{
					build.GSettings.FQueueProblem.FMainPoolCap,
					build.GSettings.FQueueProblem.FRandPoolCap,
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
	)
	return NewRawHiddenLakeNode(node, pConnsGetter, pHandlerF)
}

func NewRawHiddenLakeNode(
	pOriginNode anonymity.INode,
	pConnsGetter func() []string,
	pHandlerF handler.IHandlerF,
) IHiddenLakeNode {
	return &sHiddenLakeNode{
		fAnonNode: pOriginNode.HandleFunc(
			build.GSettings.FProtoMask.FService,
			handler.RequestHandler(pHandlerF),
		),
		fConnKeeper: connkeeper.NewConnKeeper(
			connkeeper.NewSettings(&connkeeper.SSettings{
				FDuration:    build.GSettings.GetKeeperPeriod(),
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
			uint64(build.GSettings.FProtoMask.FService),
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
			build.GSettings.FProtoMask.FService,
			pRequest.ToBytes(),
		),
	)
	if err != nil {
		return nil, err
	}
	return response.LoadResponse(rspBytes)
}

func (p *sHiddenLakeNode) GetOriginNode() anonymity.INode {
	return p.fAnonNode
}
