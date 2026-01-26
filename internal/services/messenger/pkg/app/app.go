package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/services/messenger/internal/database"
	"github.com/number571/hidden-lake/internal/services/messenger/internal/message"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/app/config"

	pkg_config "github.com/number571/hidden-lake/internal/services/messenger/pkg/config"
	hls_messenger_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/closer"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	internal_types "github.com/number571/hidden-lake/internal/utils/types"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fState state.IState

	fConfig config.IConfig
	fPathTo string

	fDatabase       database.IKVDatabase
	fIntServiceHTTP *http.Server
	fExtServiceHTTP *http.Server

	fHTTPLogger logger.ILogger
	fStdfLogger logger.ILogger
}

func NewApp(
	pCfg config.IConfig,
	pPathTo string,
) types.IRunner {
	logging := pCfg.GetLogging()

	var (
		httpLogger = std_logger.NewStdLogger(logging, http_logger.GetLogFunc())
		stdfLogger = std_logger.NewStdLogger(logging, std_logger.GetLogFunc())
	)

	return &sApp{
		fState:      state.NewBoolState(),
		fConfig:     pCfg,
		fPathTo:     pPathTo,
		fHTTPLogger: httpLogger,
		fStdfLogger: stdfLogger,
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	services := []internal_types.IServiceF{
		p.runExternalListenerHTTP,
		p.runInternalListenerHTTP,
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

		msgBroker := message.NewMessageBroker()
		hlkClient := hlk_client.NewClient(
			hlk_client.NewBuilder(),
			hlk_client.NewRequester(
				p.fConfig.GetConnection(),
				&http.Client{Timeout: build.GetSettings().GetHttpCallbackTimeout()},
			),
		)

		p.initExternalServiceHTTP(pCtx, hlkClient, msgBroker)
		p.initInternalServiceHTTP(pCtx, hlkClient, msgBroker)

		p.fStdfLogger.PushInfo(fmt.Sprintf(
			"%s is started; %s",
			hls_messenger_settings.GetAppShortNameFMT(),
			encoding.SerializeJSON(pkg_config.GetConfigSettings(p.fConfig)),
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
			hls_messenger_settings.GetAppShortNameFMT(),
		))
		return p.stop()
	}
}

func (p *sApp) runInternalListenerHTTP(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()
	defer func() { <-pCtx.Done() }()

	if p.fConfig.GetAddress().GetInternal() == "" {
		return
	}

	go func() {
		err := p.fIntServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			pChErr <- err
			return
		}
	}()
}

func (p *sApp) runExternalListenerHTTP(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()
	defer func() { <-pCtx.Done() }()

	if p.fConfig.GetAddress().GetExternal() == "" {
		return
	}

	go func() {
		err := p.fExtServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			pChErr <- err
			return
		}
	}()
}

func (p *sApp) stop() error {
	closer := closer.NewCloser(
		p.fIntServiceHTTP,
		p.fExtServiceHTTP,
		p.fDatabase,
	)
	if err := closer.Close(); err != nil {
		return errors.Join(ErrClose, err)
	}
	return nil
}
