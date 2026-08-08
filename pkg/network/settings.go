package network

import (
	"time"

	gopeer_logger "github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/pkg/network/adapters"
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
	FIsDisabled        bool
	FGeneratePeriod    time.Duration
	FNumberOfConsumers uint64
}

type SServeSettings struct {
	FLogger       gopeer_logger.ILogger
	FServiceName  string
	FFetchTimeout time.Duration
	FPowParallel  uint64
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

	if p.FQBPSettings.FGeneratePeriod == 0 {
		p.FQBPSettings.FGeneratePeriod = CDefaultQueuePeriod
	}

	if p.FQBPSettings.FNumberOfConsumers == 0 {
		p.FQBPSettings.FNumberOfConsumers = CDefaultQBPConsumers
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

	if p.FServeSettings.FFetchTimeout == 0 {
		p.FServeSettings.FFetchTimeout = CDefaultFetchTimeout
	}

	if p.FServeSettings.FPowParallel == 0 {
		p.FServeSettings.FPowParallel = CDefaultPowParallel
	}

	return p
}

func (p *sSettings) GetAdapterSettings() adapters.ISettings {
	return p.FAdapterSettings
}

func (p *sSettings) GetQBPSettings() IQBPSettings {
	return p.FQBPSettings
}

func (p *sSettings) GetFetchTimeout() time.Duration {
	return p.FServeSettings.FFetchTimeout
}

func (p *sSettings) GetPowParallel() uint64 {
	return p.FServeSettings.FPowParallel
}

func (p *sSettings) GetFmtAppName() string {
	return p.FServeSettings.FServiceName
}

func (p *sSettings) GetLogger() gopeer_logger.ILogger {
	return p.FServeSettings.FLogger
}

func (p *SQBPSettings) GetIsDisabled() bool {
	return p.FIsDisabled
}

func (p *SQBPSettings) GetGeneratePeriod() time.Duration {
	return p.FGeneratePeriod
}

func (p *SQBPSettings) GetNumberOfConsumers() uint64 {
	return p.FNumberOfConsumers
}
