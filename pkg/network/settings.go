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
	CDefaultServiceMask  = uint32(0x5f686c5f)
	CDefaultFetchTimeout = time.Minute
	CDefaultQueuePeriod  = 5 * time.Second
	CDefaultMainPoolCap  = 256
	CDefaultRandPoolCap  = 32
	CDefaultPowParallel  = 1
	CDefaultQBPConsumers = 1
)

type SSettings sSettings
type sSettings struct {
	FQBPSettings     *SQBPSettings
	FSrvSettings     *SSrvSettings
	FAdapterSettings adapters.ISettings
}

type SQBPSettings struct {
	FQueuePeriod  time.Duration
	FFetchTimeout time.Duration
	FPowParallel  uint64
	FQBPConsumers uint64
	FQueuePoolCap [2]uint64
}

type SSrvSettings struct {
	FLogger      gopeer_logger.ILogger
	FServiceMask uint32
	FServiceName string
}

func NewSettings(pSett *SSettings) ISettings {
	if pSett == nil {
		pSett = &SSettings{}
	}
	return (&sSettings{
		FAdapterSettings: pSett.FAdapterSettings,
		FQBPSettings:     pSett.FQBPSettings,
		FSrvSettings:     pSett.FSrvSettings,
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

	if p.FQBPSettings.FQueuePoolCap[0] == 0 {
		p.FQBPSettings.FQueuePoolCap[0] = CDefaultMainPoolCap
	}

	if p.FQBPSettings.FQueuePoolCap[1] == 0 {
		p.FQBPSettings.FQueuePoolCap[1] = CDefaultRandPoolCap
	}

	if p.FQBPSettings.FPowParallel == 0 {
		p.FQBPSettings.FPowParallel = CDefaultPowParallel
	}

	if p.FQBPSettings.FQBPConsumers == 0 {
		p.FQBPSettings.FQBPConsumers = CDefaultQBPConsumers
	}

	if p.FSrvSettings == nil {
		p.FSrvSettings = &SSrvSettings{}
	}

	if p.FSrvSettings.FServiceName == "" {
		p.FSrvSettings.FServiceName = CDefaultServiceName
	}

	if p.FSrvSettings.FServiceMask == 0 {
		p.FSrvSettings.FServiceMask = CDefaultServiceMask
	}

	if p.FSrvSettings.FLogger == nil {
		p.FSrvSettings.FLogger = gopeer_logger.NewLogger(
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

func (p *sSettings) GetQueuePoolCap() [2]uint64 {
	return p.FQBPSettings.FQueuePoolCap
}

func (p *sSettings) GetServiceMask() uint32 {
	return p.FSrvSettings.FServiceMask
}

func (p *sSettings) GetServiceName() string {
	return p.FSrvSettings.FServiceName
}

func (p *sSettings) GetLogger() gopeer_logger.ILogger {
	return p.FSrvSettings.FLogger
}
