package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/go-peer/pkg/utils"
	"github.com/number571/hidden-lake/cmd/applications/filesharer/internal/config"

	pkg_config "github.com/number571/hidden-lake/cmd/applications/filesharer/pkg/config"
	hlf_settings "github.com/number571/hidden-lake/cmd/applications/filesharer/pkg/settings"
	hls_client "github.com/number571/hidden-lake/cmd/service/pkg/client"
	"github.com/number571/hidden-lake/internal/closer"
	http_logger "github.com/number571/hidden-lake/internal/logger/http"
	std_logger "github.com/number571/hidden-lake/internal/logger/std"
	internal_types "github.com/number571/hidden-lake/internal/types"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fState state.IState

	fConfig  config.IConfig
	fStgPath string

	fIntServiceHTTP *http.Server
	fIncServiceHTTP *http.Server
	fServicePPROF   *http.Server

	fHTTPLogger logger.ILogger
	fStdfLogger logger.ILogger
}

func NewApp(
	pCfg config.IConfig,
	pPathTo string,
) types.IRunner {
	httpLogger := std_logger.NewStdLogger(pCfg.GetLogging(), http_logger.GetLogFunc())
	stdfLogger := std_logger.NewStdLogger(pCfg.GetLogging(), std_logger.GetLogFunc())

	return &sApp{
		fState:      state.NewBoolState(),
		fConfig:     pCfg,
		fStgPath:    filepath.Join(pPathTo, hlf_settings.CPathSTG),
		fHTTPLogger: httpLogger,
		fStdfLogger: stdfLogger,
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	services := []internal_types.IServiceF{
		p.runListenerPPROF,
		p.runIncomingListenerHTTP,
		p.runInterfaceListenerHTTP,
	}

	ctx, cancel := context.WithCancel(pCtx)
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(len(services))

	if err := p.fState.Enable(p.enable(ctx)); err != nil {
		return utils.MergeErrors(ErrRunning, err)
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
		return utils.MergeErrors(ErrService, err)
	}
}

func (p *sApp) enable(pCtx context.Context) state.IStateF {
	return func() error {
		if err := p.initStorage(); err != nil {
			return utils.MergeErrors(ErrInitSTG, err)
		}

		hlsClient := hls_client.NewClient(
			hls_client.NewBuilder(),
			hls_client.NewRequester(
				"http://"+p.fConfig.GetConnection(),
				&http.Client{Timeout: time.Hour},
			),
		)

		p.initServicePPROF()
		p.initIncomingServiceHTTP(pCtx, hlsClient)
		p.initInterfaceServiceHTTP(pCtx, hlsClient)

		p.fStdfLogger.PushInfo(fmt.Sprintf(
			"%s is started; %s",
			hlf_settings.CServiceName,
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
			hlf_settings.CServiceName,
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

func (p *sApp) runInterfaceListenerHTTP(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()

	go func() {
		err := p.fIntServiceHTTP.ListenAndServe()
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
	err := closer.CloseAll([]types.ICloser{
		p.fIntServiceHTTP,
		p.fIncServiceHTTP,
		p.fServicePPROF,
	})
	if err != nil {
		return utils.MergeErrors(ErrClose, err)
	}
	return nil
}
