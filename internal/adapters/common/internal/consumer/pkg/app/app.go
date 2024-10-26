package app

import (
	"context"
	"net/http"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/internal/adapters"
	"github.com/number571/hidden-lake/internal/adapters/common/internal/config"
	"github.com/number571/hidden-lake/internal/adapters/common/internal/consumer/internal/adapted"
	hlt_client "github.com/number571/hidden-lake/internal/helpers/traffic/pkg/client"
	"github.com/number571/hidden-lake/internal/utils/logger/std"
)

var (
	_ types.IRunner = &sApp{}
)

type sApp struct {
	fState       state.IState
	fStdfLogger  logger.ILogger
	fSettings    net_message.ISettings
	fHltAddr     string
	fServiceAddr string
	fWaitTime    time.Duration
}

func NewApp(pCfg config.IConfig) types.IRunner {
	return &sApp{
		fState:       state.NewBoolState(),
		fStdfLogger:  std.NewStdLogger(pCfg.GetLogging(), std.GetLogFunc()),
		fHltAddr:     pCfg.GetConnection().GetHLTHost(),
		fServiceAddr: pCfg.GetConnection().GetSrvHost(),
		fSettings:    pCfg.GetSettings(),
		fWaitTime:    time.Duration(pCfg.GetSettings().GetWaitTimeMS()) * time.Millisecond,
	}
}

func (p *sApp) Run(pCtx context.Context) error {
	if err := p.fState.Enable(nil); err != nil {
		return err
	}
	defer func() { _ = p.fState.Disable(nil) }()

	return adapters.ConsumeProcessor(
		pCtx,
		adapted.NewAdaptedConsumer(p.fSettings, p.fServiceAddr),
		p.fStdfLogger,
		hlt_client.NewClient(
			hlt_client.NewBuilder(),
			hlt_client.NewRequester(
				"http://"+p.fHltAddr,
				&http.Client{Timeout: time.Minute},
				p.fSettings,
			),
		),
		p.fWaitTime,
	)
}
