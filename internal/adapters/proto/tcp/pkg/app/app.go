package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/internal/adapters/proto/tcp/internal/storage"
	"github.com/number571/hidden-lake/internal/adapters/proto/tcp/pkg/app/config"
	hla_tcp_settings "github.com/number571/hidden-lake/internal/adapters/proto/tcp/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/closer"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	internal_types "github.com/number571/hidden-lake/internal/utils/types"
	"github.com/number571/hidden-lake/pkg/adapters"
	hla_http "github.com/number571/hidden-lake/pkg/adapters/http"
	hla_tcp "github.com/number571/hidden-lake/pkg/adapters/tcp"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fState   state.IState
	fWrapper config.IWrapper

	fPathTo   string
	fStorage  storage.IMessageStorage
	fDatabase database.IKVDatabase

	fHTTPLogger logger.ILogger
	fStdfLogger logger.ILogger

	fTCPAdapter    hla_tcp.ITCPAdapter
	fHTTPAdapter   hla_http.IHTTPAdapter
	fMergedAdapter adapters.IRunnerAdapter

	fServicePPROF *http.Server
}

func NewApp(pCfg config.IConfig, pPathTo string) types.IRunner {
	adaptersSettings := adapters.NewSettings(&adapters.SSettings{
		FMessageSizeBytes: pCfg.GetSettings().GetMessageSizeBytes(),
		FWorkSizeBits:     pCfg.GetSettings().GetWorkSizeBits(),
		FNetworkKey:       pCfg.GetSettings().GetNetworkKey(),
	})
	tcpAdapter := hla_tcp.NewTCPAdapter(
		hla_tcp.NewSettings(&hla_tcp.SSettings{
			FAddress:         pCfg.GetAddress().GetExternal(),
			FAdapterSettings: adaptersSettings,
		}),
		func() []string { return pCfg.GetConnections() },
	)
	httpAdapter := hla_http.NewHTTPAdapter(
		hla_http.NewSettings(&hla_http.SSettings{
			FAddress:         pCfg.GetAddress().GetInternal(),
			FAdapterSettings: adaptersSettings,
		}),
		func() []string { return pCfg.GetEndpoints() },
	)
	return &sApp{
		fState:         state.NewBoolState(),
		fPathTo:        pPathTo,
		fWrapper:       config.NewWrapper(pCfg),
		fStdfLogger:    std_logger.NewStdLogger(pCfg.GetLogging(), std_logger.GetLogFunc()),
		fHTTPLogger:    std_logger.NewStdLogger(pCfg.GetLogging(), http_logger.GetLogFunc()),
		fTCPAdapter:    tcpAdapter,
		fHTTPAdapter:   httpAdapter,
		fMergedAdapter: adapters.NewMergedRunnerAdapter(tcpAdapter, httpAdapter),
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	services := []internal_types.IServiceF{
		p.runMergedAdapter,
		p.runAdaptedRelayer,
		p.runListenerPPROF,
	}

	ctx, cancel := context.WithCancel(pCtx)
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(len(services))

	if err := p.fState.Enable(p.enable(ctx)); err != nil {
		return errors.Join(ErrRunning, err)
	}
	defer func() { _ = p.fState.Disable(p.disable(cancel, wg)) }()

	chErr := make(chan error, len(services))
	for _, f := range services {
		go f(ctx, wg, chErr)
	}

	select {
	case <-pCtx.Done():
		return pCtx.Err()
	case err := <-chErr:
		return errors.Join(ErrService, err)
	}
}

func (p *sApp) enable(pCtx context.Context) state.IStateF {
	return func() error {
		if err := p.initDatabase(); err != nil {
			return errors.Join(ErrInitDB, err)
		}
		p.initStorage(p.fDatabase)

		p.initHandlers(pCtx)
		p.initServicePPROF()

		p.fStdfLogger.PushInfo(fmt.Sprintf( // nolint: perfsprint
			"%s is started",
			hla_tcp_settings.GServiceName.Short(),
		))
		return nil
	}
}

func (p *sApp) disable(pCancel context.CancelFunc, pWg *sync.WaitGroup) state.IStateF {
	return func() error {
		pCancel()
		pWg.Wait() // wait canceled context

		p.fStdfLogger.PushInfo(fmt.Sprintf( // nolint: perfsprint
			"%s is stopped",
			hla_tcp_settings.GServiceName.Short(),
		))
		return p.stop()
	}
}

func (p *sApp) runMergedAdapter(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()

	if err := p.fMergedAdapter.Run(pCtx); err != nil {
		pChErr <- err
		return
	}
}

func (p *sApp) runAdaptedRelayer(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()

	for {
		select {
		case <-pCtx.Done():
			pChErr <- pCtx.Err()
			return
		default:
			msg, err := p.fMergedAdapter.Consume(pCtx)
			if err != nil {
				continue
			}

			if err := p.fStorage.Push(msg); err != nil {
				continue
			}
			_ = p.fMergedAdapter.Produce(pCtx, msg)
		}
	}
}

func (p *sApp) runListenerPPROF(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()

	if p.fWrapper.GetConfig().GetAddress().GetPPROF() == "" {
		return
	}

	go func() {
		err := p.fServicePPROF.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			pChErr <- err
			return
		}
	}()

	<-pCtx.Done()
}

func (p *sApp) stop() error {
	err := closer.CloseAll([]io.Closer{
		p.fDatabase,
		p.fServicePPROF,
	})
	if err != nil {
		return errors.Join(ErrClose, err)
	}
	return nil
}
