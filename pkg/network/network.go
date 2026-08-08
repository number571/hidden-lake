package network

import (
	"context"
	"errors"
	"sync"

	anonymity "github.com/number571/go-peer/pkg/anonymity/qb"
	"github.com/number571/go-peer/pkg/anonymity/qb/queue"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/pkg/network/adapters"
	"github.com/number571/hidden-lake/pkg/network/handler"
	"github.com/number571/hidden-lake/pkg/network/request"
	"github.com/number571/hidden-lake/pkg/network/response"
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
	pScheme layer2.IScheme,
	pKeysContainer layer2.IKeysContainer,
	pKVDatabase database.IKVDatabase,
	pRunnerAdapter adapters.IRunnerAdapter,
	pHandlerF handler.IHandlerF,
) (IHiddenLakeNode, error) {
	if pScheme.GetPayloadLimit() <= anonymity.CMessageHeadSize {
		return nil, ErrPayloadLimit
	}
	return &sHiddenLakeNode{
		fSettings: pSettings,
		fOriginNode: anonymity.NewNode(
			anonymity.NewSettings(&anonymity.SSettings{
				FServiceName:  pSettings.GetFmtAppName(),
				FFetchTimeout: pSettings.GetFetchTimeout(),
			}),
			handler.RequestHandler(pHandlerF),
			pSettings.GetLogger(),
			pRunnerAdapter,
			pKVDatabase,
			pKeysContainer,
			getQueueProcessor(pSettings, pScheme),
		),
	}, nil
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
		ra := p.fOriginNode.GetAdapter().(adapters.IRunnerAdapter)
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
	pKey layer2.IParticipantKey,
	pRequest request.IRequest,
) error {
	err := p.fOriginNode.SendPayload(
		pCtx,
		pKey,
		pRequest.ToBytes(),
	)
	if err != nil {
		return errors.Join(ErrSendRequest, err)
	}
	return nil
}

func (p *sHiddenLakeNode) FetchRequest(
	pCtx context.Context,
	pKey layer2.IParticipantKey,
	pRequest request.IRequest,
) (response.IResponse, error) {
	rspBytes, err := p.fOriginNode.FetchPayload(
		pCtx,
		pKey,
		pRequest.ToBytes(),
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

func getQueueProcessor(
	pSettings ISettings,
	pScheme layer2.IScheme,
) queue.IQBProblemProcessor {
	buildSettings := build.GetSettings()
	qbpSettings := pSettings.GetQBPSettings()

	qSettings := queue.NewSettings(&queue.SSettings{
		FMessageConstructSettings: layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: pSettings.GetAdapterSettings(),
			FParallel: pSettings.GetPowParallel(),
		}),
		FQueuePeriod:  qbpSettings.GetGeneratePeriod(),
		FConsumersCap: qbpSettings.GetNumberOfConsumers(),
		FQueuePoolCap: buildSettings.FStorageManager.FQueuePoolCap,
	})

	processor := queue.NewQBProblemProcessor
	if qbpSettings.GetIsDisabled() {
		processor = newSilentQueueProcessor
	}

	return processor(qSettings, pScheme)
}
