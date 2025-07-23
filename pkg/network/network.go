package network

import (
	"context"
	"errors"
	"sync"

	"github.com/number571/go-peer/pkg/anonymity"
	"github.com/number571/go-peer/pkg/anonymity/queue"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/pkg/adapters"
	"github.com/number571/hidden-lake/pkg/handler"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"
)

var (
	_ IHiddenLakeNode = &sHiddenLakeNode{}
)

type sHiddenLakeNode struct {
	fSettings   ISettings
	fOriginNode anonymity.INode
}

func NewHiddenLakeNode(
	pSettings ISettings,
	pPrivKey asymmetric.IPrivKey,
	pKVDatabase database.IKVDatabase,
	pRunnerAdapter adapters.IRunnerAdapter,
	pHandlerF handler.IHandlerF,
) IHiddenLakeNode {
	buildSettings := build.GetSettings()
	adaptersSettings := pSettings.GetAdapterSettings()
	return &sHiddenLakeNode{
		fSettings: pSettings,
		fOriginNode: anonymity.NewNode(
			anonymity.NewSettings(&anonymity.SSettings{
				FServiceName:  pSettings.GetFmtAppName(),
				FFetchTimeout: pSettings.GetFetchTimeout(),
			}),
			pSettings.GetLogger(),
			pRunnerAdapter,
			pKVDatabase,
			queue.NewQBProblemProcessor(
				queue.NewSettings(&queue.SSettings{
					FMessageConstructSettings: layer1.NewConstructSettings(&layer1.SConstructSettings{
						FSettings: adaptersSettings,
						FParallel: pSettings.GetPowParallel(),
					}),
					FNetworkMask:  buildSettings.FProtoMask.FNetwork,
					FQueuePeriod:  pSettings.GetQueuePeriod(),
					FConsumersCap: pSettings.GetQBPConsumers(),
					FQueuePoolCap: buildSettings.FStorageManager.FQueuePoolCap,
				}),
				func() client.IClient {
					client := client.NewClient(pPrivKey, adaptersSettings.GetMessageSizeBytes())
					if client.GetPayloadLimit() <= encoding.CSizeUint64 {
						panic(`client.GetPayloadLimit() <= encoding.CSizeUint64`)
					}
					return client
				}(),
			),
		).HandleFunc(
			buildSettings.FProtoMask.FService,
			handler.RequestHandler(pHandlerF),
		),
	}
}

func (p *sHiddenLakeNode) GetOriginNode() anonymity.INode {
	return p.fOriginNode
}

func (p *sHiddenLakeNode) Run(pCtx context.Context) error {
	chCtx, cancel := context.WithCancel(pCtx)
	defer cancel()

	const N = 2

	errs := make([]error, N)
	wg := &sync.WaitGroup{}
	wg.Add(N)

	go func() {
		defer func() { wg.Done(); cancel() }()
		ra, ok := p.fOriginNode.GetAdapter().(adapters.IRunnerAdapter)
		if !ok {
			errs[0] = ErrAdapterNotRunner
			return
		}
		errs[0] = ra.Run(chCtx)
	}()

	go func() {
		defer func() { wg.Done(); cancel() }()
		errs[1] = p.fOriginNode.Run(chCtx)
	}()

	wg.Wait()

	select {
	case <-pCtx.Done():
		return pCtx.Err()
	default:
		errs := append([]error{ErrRunning}, errs...)
		return errors.Join(errs...)
	}
}

func (p *sHiddenLakeNode) SendRequest(
	pCtx context.Context,
	pPubKey asymmetric.IPubKey,
	pRequest request.IRequest,
) error {
	buildSettings := build.GetSettings()
	err := p.fOriginNode.SendPayload(
		pCtx,
		pPubKey,
		payload.NewPayload64(
			uint64(buildSettings.FProtoMask.FService),
			pRequest.ToBytes(),
		),
	)
	if err != nil {
		return errors.Join(ErrSendRequest, err)
	}
	return nil
}

func (p *sHiddenLakeNode) FetchRequest(
	pCtx context.Context,
	pPubKey asymmetric.IPubKey,
	pRequest request.IRequest,
) (response.IResponse, error) {
	buildSettings := build.GetSettings()
	rspBytes, err := p.fOriginNode.FetchPayload(
		pCtx,
		pPubKey,
		payload.NewPayload32(
			buildSettings.FProtoMask.FService,
			pRequest.ToBytes(),
		),
	)
	if err != nil {
		return nil, errors.Join(ErrFetchRequest, err)
	}
	rsp, err := response.LoadResponse(rspBytes)
	if err != nil {
		return nil, errors.Join(ErrLoadResponse, err)
	}
	return rsp, nil
}
