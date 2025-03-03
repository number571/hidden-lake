package app

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/internal/composite/pkg/app/config"

	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	internal_types "github.com/number571/hidden-lake/internal/utils/types"

	hlc_settings "github.com/number571/hidden-lake/internal/composite/pkg/settings"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fState   state.IState
	fConfig  config.IConfig
	fRunners []types.IRunner

	fStdfLogger logger.ILogger
}

func NewApp(
	pCfg config.IConfig,
	pRunners []types.IRunner,
) types.IRunner {
	stdfLogger := std_logger.NewStdLogger(pCfg.GetLogging(), std_logger.GetLogFunc())

	return &sApp{
		fState:      state.NewBoolState(),
		fConfig:     pCfg,
		fRunners:    pRunners,
		fStdfLogger: stdfLogger,
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	services := runnersToServices(p.fRunners)

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

func (p *sApp) enable(_ context.Context) state.IStateF {
	return func() error {
		p.fStdfLogger.PushInfo(fmt.Sprintf(
			"%s is started; %s",
			hlc_settings.GServiceName.Short(),
			encoding.SerializeJSON(p.fConfig.GetServices()),
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
			hlc_settings.GServiceName.Short(),
		))
		return nil
	}
}

func runnersToServices(pRunners []types.IRunner) []internal_types.IServiceF {
	services := make([]internal_types.IServiceF, 0, len(pRunners))
	for _, runner := range pRunners {
		runner := runner
		services = append(
			services,
			func(pCtx context.Context, pWg *sync.WaitGroup, pChErr chan<- error) {
				defer pWg.Done()
				if err := runner.Run(pCtx); err != nil {
					pChErr <- err
					return
				}
			},
		)
	}
	return services
}
