package app

import (
	"context"

	"github.com/number571/go-peer/pkg/logger"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/internal/adapters"
	"github.com/number571/hidden-lake/internal/adapters/common/internal/config"
	"github.com/number571/hidden-lake/internal/adapters/common/internal/producer/internal/adapted"
	"github.com/number571/hidden-lake/internal/utils/logger/std"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fState        state.IState
	fStdfLogger   logger.ILogger
	fSettings     net_message.ISettings
	fServiceAddr  string
	fIncomingAddr string
}

func NewApp(pCfg config.IConfig) types.IRunner {
	return &sApp{
		fState:        state.NewBoolState(),
		fStdfLogger:   std.NewStdLogger(pCfg.GetLogging(), std.GetLogFunc()),
		fServiceAddr:  pCfg.GetConnection().GetSrvHost(),
		fIncomingAddr: pCfg.GetAddress(),
		fSettings:     pCfg.GetSettings(),
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	if err := p.fState.Enable(nil); err != nil {
		return err
	}
	defer func() { _ = p.fState.Disable(nil) }()

	return adapters.ProduceProcessor(
		pCtx,
		adapted.NewAdaptedProducer(p.fServiceAddr),
		p.fStdfLogger,
		p.fSettings,
		p.fIncomingAddr,
	)
}
