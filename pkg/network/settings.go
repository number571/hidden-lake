package network

import (
	"time"

	gopeer_logger "github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/pkg/adapters"
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FSubSettings     *SSubSettings
	FQueuePeriod     time.Duration
	FFetchTimeout    time.Duration
	FAdapterSettings adapters.ISettings
}

type SSubSettings struct {
	FLogger       gopeer_logger.ILogger
	FPowParallel  uint64
	FQBPConsumers uint64
	FServiceName  string
}

func NewSettings(pSett *SSettings) ISettings {
	if pSett == nil {
		pSett = &SSettings{}
	}
	return (&sSettings{
		FAdapterSettings: pSett.FAdapterSettings,
		FQueuePeriod:     pSett.FQueuePeriod,
		FFetchTimeout:    pSett.FFetchTimeout,
		FSubSettings:     pSett.FSubSettings,
	}).useDefault()
}

func NewSettingsByNetworkKey(pNetworkKey string, pSubSettings *SSubSettings) ISettings {
	network, ok := build.GetNetwork(pNetworkKey)
	if !ok {
		panic("network not found")
	}
	return NewSettings(&SSettings{
		FAdapterSettings: adapters.NewSettings(&adapters.SSettings{
			FNetworkKey:       pNetworkKey,
			FWorkSizeBits:     network.FWorkSizeBits,
			FMessageSizeBytes: network.FMessageSizeBytes,
		}),
		FQueuePeriod:  network.GetQueuePeriod(),
		FFetchTimeout: network.GetFetchTimeout(),
		FSubSettings:  pSubSettings,
	})
}

func (p *sSettings) useDefault() *sSettings {
	defaultNetwork, _ := build.GetNetwork(build.CDefaultNetwork)

	if p.FAdapterSettings == nil {
		p.FAdapterSettings = adapters.NewSettings(nil)
	}

	if p.FQueuePeriod == 0 {
		p.FQueuePeriod = defaultNetwork.GetQueuePeriod()
	}

	if p.FFetchTimeout == 0 {
		p.FFetchTimeout = defaultNetwork.GetFetchTimeout()
	}

	if p.FSubSettings == nil {
		p.FSubSettings = &SSubSettings{}
	}

	if p.FSubSettings.FServiceName == "" {
		p.FSubSettings.FServiceName = "_"
	}

	if p.FSubSettings.FPowParallel == 0 {
		p.FSubSettings.FPowParallel = 1
	}

	if p.FSubSettings.FQBPConsumers == 0 {
		p.FSubSettings.FQBPConsumers = 1
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
	return p.FQueuePeriod
}

func (p *sSettings) GetFetchTimeout() time.Duration {
	return p.FFetchTimeout
}

func (p *sSettings) GetPowParallel() uint64 {
	return p.FSubSettings.FPowParallel
}

func (p *sSettings) GetQBPConsumers() uint64 {
	return p.FSubSettings.FQBPConsumers
}

func (p *sSettings) GetServiceName() string {
	return p.FSubSettings.FServiceName
}

func (p *sSettings) GetLogger() gopeer_logger.ILogger {
	return p.FSubSettings.FLogger
}
