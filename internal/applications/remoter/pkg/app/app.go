package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/internal/applications/remoter/pkg/app/config"
	pkg_config "github.com/number571/hidden-lake/internal/applications/remoter/pkg/config"
	hlr_settings "github.com/number571/hidden-lake/internal/applications/remoter/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/closer"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	internal_types "github.com/number571/hidden-lake/internal/utils/types"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fState state.IState

	fConfig config.IConfig

	fHTTPLogger logger.ILogger
	fStdfLogger logger.ILogger

	fIncServiceHTTP *http.Server
	fServicePPROF   *http.Server
}

func NewApp(
	pCfg config.IConfig,
) types.IRunner {
	logging := pCfg.GetLogging()

	var (
		httpLogger = std_logger.NewStdLogger(logging, http_logger.GetLogFunc())
		stdfLogger = std_logger.NewStdLogger(logging, std_logger.GetLogFunc())
	)

	return &sApp{
		fState:      state.NewBoolState(),
		fConfig:     pCfg,
		fHTTPLogger: httpLogger,
		fStdfLogger: stdfLogger,
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	services := []internal_types.IServiceF{
		p.runIncomingListenerHTTP,
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
		p.initIncomingServiceHTTP(pCtx)
		p.initServicePPROF()

		p.fStdfLogger.PushInfo(fmt.Sprintf(
			"%s is started; %s",
			hlr_settings.GServiceName.Short(),
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
			hlr_settings.GServiceName.Short(),
		))
		return p.stop()
	}
}

func (p *sApp) runListenerPPROF(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()

	if p.fConfig.GetAddress().GetPPROF() == "" {
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

func (p *sApp) runIncomingListenerHTTP(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()

	go func() {
		err := p.fIncServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			pChErr <- err
			return
		}
	}()

	<-pCtx.Done()
}

func (p *sApp) stop() error {
	err := closer.CloseAll([]io.Closer{
		p.fIncServiceHTTP,
		p.fServicePPROF,
	})
	if err != nil {
		return errors.Join(ErrClose, err)
	}
	return nil
}
