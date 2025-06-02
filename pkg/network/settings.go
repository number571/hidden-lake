package network

import (
	"time"

	gopeer_logger "github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/pkg/adapters"
)

var (
	_ ISettings = &sSettings{}
)

const (
	CDefaultFetchTimeout = time.Minute
	CDefaultQueuePeriod  = 5 * time.Second
	CDefaultQBPConsumers = 1
	CDefaultPowParallel  = 1
)

type SSettings sSettings
type sSettings struct {
	FQBPSettings     *SQBPSettings
	FSubSettings     *SSubSettings
	FAdapterSettings adapters.ISettings
}

type SQBPSettings struct {
	FQueuePeriod  time.Duration
	FFetchTimeout time.Duration
	FPowParallel  uint64
	FQBPConsumers uint64
}

type SSubSettings struct {
	FLogger      gopeer_logger.ILogger
	FServiceName string
}

func NewSettings(pSett *SSettings) ISettings {
	if pSett == nil {
		pSett = &SSettings{}
	}
	return (&sSettings{
		FAdapterSettings: pSett.FAdapterSettings,
		FQBPSettings:     pSett.FQBPSettings,
		FSubSettings:     pSett.FSubSettings,
	}).useDefault()
}

func (p *sSettings) useDefault() *sSettings {
	if p.FAdapterSettings == nil {
		p.FAdapterSettings = adapters.NewSettings(nil)
	}

	if p.FQBPSettings == nil {
		p.FQBPSettings = &SQBPSettings{}
	}

	if p.FQBPSettings.FQueuePeriod == 0 {
		p.FQBPSettings.FQueuePeriod = CDefaultQueuePeriod
	}

	if p.FQBPSettings.FFetchTimeout == 0 {
		p.FQBPSettings.FFetchTimeout = CDefaultFetchTimeout
	}

	if p.FQBPSettings.FPowParallel == 0 {
		p.FQBPSettings.FPowParallel = CDefaultPowParallel
	}

	if p.FQBPSettings.FQBPConsumers == 0 {
		p.FQBPSettings.FQBPConsumers = CDefaultQBPConsumers
	}

	if p.FSubSettings == nil {
		p.FSubSettings = &SSubSettings{}
	}

	if p.FSubSettings.FServiceName == "" {
		p.FSubSettings.FServiceName = "_"
	}

	if p.FSubSettings.FLogger == nil {
		p.FSubSettings.FLogger = gopeer_logger.NewLogger(
			gopeer_logger.NewSettings(&gopeer_logger.SSettings{}),
			func(_ gopeer_logger.ILogArg) string { return "" },
		)
	}

	return p
}

func (p *sSettings) GetAdapterSettings() adapters.ISettings {
	return p.FAdapterSettings
}

func (p *sSettings) GetQueuePeriod() time.Duration {
	return p.FQBPSettings.FQueuePeriod
}

func (p *sSettings) GetFetchTimeout() time.Duration {
	return p.FQBPSettings.FFetchTimeout
}

func (p *sSettings) GetPowParallel() uint64 {
	return p.FQBPSettings.FPowParallel
}

func (p *sSettings) GetQBPConsumers() uint64 {
	return p.FQBPSettings.FQBPConsumers
}

func (p *sSettings) GetServiceName() string {
	return p.FSubSettings.FServiceName
}

func (p *sSettings) GetLogger() gopeer_logger.ILogger {
	return p.FSubSettings.FLogger
}
