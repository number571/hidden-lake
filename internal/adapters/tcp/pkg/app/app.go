package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/internal/adapters/tcp/pkg/app/config"
	"github.com/number571/hidden-lake/pkg/adapters/tcp"

	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/hidden-lake/internal/adapters/tcp/internal/storage"
	hla_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/closer"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	internal_types "github.com/number571/hidden-lake/internal/utils/types"
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

	fTCPAdapter   tcp.ITCPAdapter
	fServiceHTTP  *http.Server
	fServicePPROF *http.Server
}

func NewApp(
	pCfg config.IConfig,
	pPathTo string,
) types.IRunner {
	return &sApp{
		fState:      state.NewBoolState(),
		fPathTo:     pPathTo,
		fWrapper:    config.NewWrapper(pCfg),
		fStdfLogger: std_logger.NewStdLogger(pCfg.GetLogging(), std_logger.GetLogFunc()),
		fHTTPLogger: std_logger.NewStdLogger(pCfg.GetLogging(), http_logger.GetLogFunc()),
		fTCPAdapter: tcp.NewTCPAdapter(
			tcp.NewSettings(&tcp.SSettings{
				FAddress:          pCfg.GetAddress().GetExternal(),
				FMessageSizeBytes: pCfg.GetSettings().GetMessageSizeBytes(),
				FWorkSizeBits:     pCfg.GetSettings().GetWorkSizeBits(),
				FNetworkKey:       pCfg.GetSettings().GetNetworkKey(),
			}),
			func() []string { return pCfg.GetConnections() },
		),
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	services := []internal_types.IServiceF{
		p.runTCPAdapter,
		p.runAdaptedRelayer,
		p.runListenerHTTP,
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

		p.initServiceHTTP(pCtx)
		p.initServicePPROF()

		p.fStdfLogger.PushInfo(fmt.Sprintf( // nolint: perfsprint
			"%s is started",
			hla_settings.GServiceName.Short(),
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
			hla_settings.GServiceName.Short(),
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

func (p *sApp) runAdaptedRelayer(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()

	for {
		select {
		case <-pCtx.Done():
			pChErr <- pCtx.Err()
			return
		default:
			msg, err := p.fTCPAdapter.Consume(pCtx)
			if err != nil {
				continue
			}

			if err := p.fStorage.Push(msg); err != nil {
				continue
			}
			_ = p.fTCPAdapter.Produce(pCtx, msg)

			endpoints := p.fWrapper.GetConfig().GetEndpoints()

			wg := &sync.WaitGroup{}
			wg.Add(len(endpoints))

			for _, ep := range endpoints {
				ep := ep
				go func() {
					defer wg.Done()
					produceToEndpoint(pCtx, ep, msg)
				}()
			}

			wg.Wait()
		}
	}
}

func produceToEndpoint(pCtx context.Context, pEndpoint string, pNetMsg net_message.IMessage) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(
		pCtx,
		http.MethodPost,
		"http://"+pEndpoint+hla_settings.CHandleNetworkAdapterPath,
		strings.NewReader(pNetMsg.ToString()),
	)
	if err != nil {
		return
	}
	rsp, err := client.Do(req)
	if err != nil {
		return
	}
	rsp.Body.Close()
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

func (p *sApp) runListenerHTTP(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()

	if p.fWrapper.GetConfig().GetAddress().GetInternal() == "" {
		return
	}

	go func() {
		err := p.fServiceHTTP.ListenAndServe()
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
		p.fServiceHTTP,
		p.fServicePPROF,
	})
	if err != nil {
		return errors.Join(ErrClose, err)
	}
	return nil
}
