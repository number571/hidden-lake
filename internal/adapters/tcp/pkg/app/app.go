package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/adapters/tcp/pkg/app/config"
	hla_tcp_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/closer"
	anon_logger "github.com/number571/hidden-lake/internal/utils/logger/anon"
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
	fDatabase database.IKVDatabase

	fAnonLogger logger.ILogger
	fHTTPLogger logger.ILogger
	fStdfLogger logger.ILogger

	fTCPAdapter  hla_tcp.ITCPAdapter
	fHTTPAdapter hla_http.IHTTPAdapter
}

func NewApp(pCfg config.IConfig, pPathTo string) types.IRunner {
	logging := pCfg.GetLogging()
	lruCache := cache.NewLRUCache(build.GSettings.FNetworkManager.FCacheHashesCap)
	adaptersSettings := adapters.NewSettings(&adapters.SSettings{
		FMessageSizeBytes: pCfg.GetSettings().GetMessageSizeBytes(),
		FWorkSizeBits:     pCfg.GetSettings().GetWorkSizeBits(),
		FNetworkKey:       pCfg.GetSettings().GetNetworkKey(),
	})
	return &sApp{
		fState:      state.NewBoolState(),
		fPathTo:     pPathTo,
		fWrapper:    config.NewWrapper(pCfg),
		fAnonLogger: std_logger.NewStdLogger(logging, anon_logger.GetLogFunc()),
		fStdfLogger: std_logger.NewStdLogger(logging, std_logger.GetLogFunc()),
		fHTTPLogger: std_logger.NewStdLogger(logging, http_logger.GetLogFunc()),
		fTCPAdapter: hla_tcp.NewTCPAdapter(
			hla_tcp.NewSettings(&hla_tcp.SSettings{
				FAddress:         pCfg.GetAddress().GetExternal(),
				FAdapterSettings: adaptersSettings,
			}),
			lruCache,
			func() []string { return pCfg.GetConnections() },
		),
		fHTTPAdapter: hla_http.NewHTTPAdapter(
			hla_http.NewSettings(&hla_http.SSettings{
				FAddress:         pCfg.GetAddress().GetInternal(),
				FAdapterSettings: adaptersSettings,
			}),
			lruCache,
			func() []string { return pCfg.GetEndpoints() },
		),
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	services := []internal_types.IServiceF{
		p.runTCPAdapter,
		p.runTCPRelayer,
		p.runHTTPAdapter,
		p.runHTTPRelayer,
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

		p.initLoggers()
		p.initHandlers(pCtx)

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

func (p *sApp) runTCPAdapter(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()

	if err := p.fTCPAdapter.Run(pCtx); err != nil {
		pChErr <- err
		return
	}
}

func (p *sApp) runHTTPAdapter(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()

	if err := p.fHTTPAdapter.Run(pCtx); err != nil {
		pChErr <- err
		return
	}
}

func (p *sApp) runTCPRelayer(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()

	for {
		select {
		case <-pCtx.Done():
			pChErr <- pCtx.Err()
			return
		default:
			// TCP (connections) -> HTTP (endpoints), TCP (connections)
			msg, err := p.fTCPAdapter.Consume(pCtx)
			if err != nil {
				continue
			}
			if err := p.setIntoDB(msg); err != nil {
				continue
			}
			if err := p.fHTTPAdapter.Produce(pCtx, msg); err != nil {
				if !errors.Is(err, hla_http.ErrNoConnections) {
					continue
				}
			}
			_ = p.fTCPAdapter.Produce(pCtx, msg)
		}
	}
}

func (p *sApp) runHTTPRelayer(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()

	for {
		select {
		case <-pCtx.Done():
			pChErr <- pCtx.Err()
			return
		default:
			// HTTP (endpoints) -> TCP (connections)
			msg, err := p.fHTTPAdapter.Consume(pCtx)
			if err != nil {
				continue
			}
			if err := p.setIntoDB(msg); err != nil {
				continue
			}
			_ = p.fTCPAdapter.Produce(pCtx, msg)
		}
	}
}

func (p *sApp) setIntoDB(msg layer1.IMessage) error {
	_, err := p.fDatabase.Get(msg.GetHash())
	if err == nil {
		return ErrExist
	}
	if !errors.Is(err, database.ErrNotFound) {
		return err
	}
	return p.fDatabase.Set(msg.GetHash(), []byte{})
}

func (p *sApp) stop() error {
	err := closer.CloseAll([]io.Closer{
		p.fDatabase,
	})
	if err != nil {
		return errors.Join(ErrClose, err)
	}
	return nil
}
