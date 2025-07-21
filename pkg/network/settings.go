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
	CDefaultServiceName  = "_"
	CDefaultFetchTimeout = time.Minute
	CDefaultQueuePeriod  = 5 * time.Second
	CDefaultPowParallel  = 1
	CDefaultQBPConsumers = 1
)

type SSettings sSettings
type sSettings struct {
	FQBPSettings     *SQBPSettings
	FServeSettings   *SServeSettings
	FAdapterSettings adapters.ISettings
}

type SQBPSettings struct {
	FQueuePeriod  time.Duration
	FFetchTimeout time.Duration
	FPowParallel  uint64
	FQBPConsumers uint64
}

type SServeSettings struct {
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
		FServeSettings:   pSett.FServeSettings,
	}).initDefault()
}

func (p *sSettings) initDefault() *sSettings {
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

	if p.FServeSettings == nil {
		p.FServeSettings = &SServeSettings{}
	}

	if p.FServeSettings.FServiceName == "" {
		p.FServeSettings.FServiceName = CDefaultServiceName
	}

	if p.FServeSettings.FLogger == nil {
		p.FServeSettings.FLogger = gopeer_logger.NewLogger(
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

func (p *sSettings) GetAppName() string {
	return p.FServeSettings.FServiceName
}

func (p *sSettings) GetLogger() gopeer_logger.ILogger {
	return p.FServeSettings.FLogger
}
