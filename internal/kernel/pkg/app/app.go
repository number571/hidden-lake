package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/internal/kernel/pkg/app/config"
	"github.com/number571/hidden-lake/internal/utils/closer"
	"github.com/number571/hidden-lake/pkg/network"

	pkg_config "github.com/number571/hidden-lake/internal/kernel/pkg/config"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	anon_logger "github.com/number571/hidden-lake/internal/utils/logger/anon"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	internal_types "github.com/number571/hidden-lake/internal/utils/types"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fState  state.IState
	fPathTo string

	fCfgW    config.IWrapper
	fNode    network.IHiddenLakeNode
	fPrivKey asymmetric.IPrivKey

	fAnonLogger logger.ILogger
	fHTTPLogger logger.ILogger
	fStdfLogger logger.ILogger

	fServiceHTTP *http.Server
}

func NewApp(
	pCfg config.IConfig,
	pPrivKey asymmetric.IPrivKey,
	pPathTo string,
) types.IRunner {
	logging := pCfg.GetLogging()

	var (
		anonLogger = std_logger.NewStdLogger(logging, anon_logger.GetLogFunc())
		httpLogger = std_logger.NewStdLogger(logging, http_logger.GetLogFunc())
		stdfLogger = std_logger.NewStdLogger(logging, std_logger.GetLogFunc())
	)

	return &sApp{
		fState:      state.NewBoolState(),
		fPathTo:     pPathTo,
		fCfgW:       config.NewWrapper(pCfg),
		fPrivKey:    pPrivKey,
		fAnonLogger: anonLogger,
		fHTTPLogger: httpLogger,
		fStdfLogger: stdfLogger,
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	services := []internal_types.IServiceF{
		p.runListenerInternal,
		p.runAnonymityNode,
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
		if err := p.initAnonNode(); err != nil {
			return errors.Join(ErrCreateAnonNode, err)
		}

		p.initServiceHTTP(pCtx)

		p.fStdfLogger.PushInfo(fmt.Sprintf(
			"%s is started; %s",
			hlk_settings.GetServiceName().Short(),
			encoding.SerializeJSON(pkg_config.GetConfigSettings(
				p.fCfgW.GetConfig(),
				p.fNode.GetOriginNode().GetQBProcessor().GetClient(),
			)),
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
			hlk_settings.GetServiceName().Short(),
		))
		return p.stop()
	}
}

func (p *sApp) stop() error {
	err := closer.CloseAll([]io.Closer{
		p.fServiceHTTP,
		p.fNode.GetOriginNode().GetKVDatabase(),
	})
	if err != nil {
		return errors.Join(ErrClose, err)
	}
	return nil
}

func (p *sApp) runAnonymityNode(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()

	if err := p.fNode.Run(pCtx); err != nil {
		pChErr <- err
		return
	}
}

func (p *sApp) runListenerInternal(pCtx context.Context, wg *sync.WaitGroup, pChErr chan<- error) {
	defer wg.Done()
	defer func() { <-pCtx.Done() }()

	if p.fCfgW.GetConfig().GetAddress().GetInternal() == "" {
		return
	}

	go func() {
		err := p.fServiceHTTP.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			pChErr <- err
			return
		}
	}()
}
