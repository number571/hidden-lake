package network

import (
	"time"

	gopeer_logger "github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/build"
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FSubSettings      *SSubSettings
	FQueuePeriod      time.Duration
	FFetchTimeout     time.Duration
	FMessageSizeBytes uint64
	FWorkSizeBits     uint64
	FNetworkKey       string
}

type SSubSettings struct {
	FLogger      gopeer_logger.ILogger
	FParallel    uint64
	FTCPAddress  string
	FServiceName string
}

func NewSettings(pSett *SSettings) ISettings {
	if pSett == nil {
		pSett = &SSettings{}
	}
	return (&sSettings{
		FMessageSizeBytes: pSett.FMessageSizeBytes,
		FQueuePeriod:      pSett.FQueuePeriod,
		FFetchTimeout:     pSett.FFetchTimeout,
		FSubSettings:      pSett.FSubSettings,
		FWorkSizeBits:     pSett.FWorkSizeBits,
		FNetworkKey:       pSett.FNetworkKey,
	}).useDefault()
}

func NewSettingsByNetworkKey(pNetworkKey string, pSubSettings *SSubSettings) ISettings {
	network, ok := build.GNetworks[pNetworkKey]
	if !ok {
		panic("network not found")
	}
	return NewSettings(&SSettings{
		FWorkSizeBits:     network.FWorkSizeBits,
		FNetworkKey:       pNetworkKey,
		FMessageSizeBytes: network.FMessageSizeBytes,
		FQueuePeriod:      network.GetQueuePeriod(),
		FFetchTimeout:     network.GetFetchTimeout(),
		FSubSettings:      pSubSettings,
	})
}

func (p *sSettings) useDefault() *sSettings {
	defaultNetwork := build.GNetworks[build.CDefaultNetwork]

	if p.FMessageSizeBytes == 0 {
		p.FMessageSizeBytes = defaultNetwork.FMessageSizeBytes
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

	if p.FSubSettings.FLogger == nil {
		p.FSubSettings.FLogger = gopeer_logger.NewLogger(
			gopeer_logger.NewSettings(&gopeer_logger.SSettings{}),
			func(_ gopeer_logger.ILogArg) string { return "" },
		)
	}

	return p
}

func (p *sSettings) GetNetworkKey() string {
	return p.FNetworkKey
}

func (p *sSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}

func (p *sSettings) GetMessageSizeBytes() uint64 {
	return p.FMessageSizeBytes
}

func (p *sSettings) GetQueuePeriod() time.Duration {
	return p.FQueuePeriod
}

func (p *sSettings) GetFetchTimeout() time.Duration {
	return p.FFetchTimeout
}

func (p *sSettings) GetParallel() uint64 {
	return p.FSubSettings.FParallel
}

func (p *sSettings) GetServiceName() string {
	return p.FSubSettings.FServiceName
}

func (p *sSettings) GetTCPAddress() string {
	return p.FSubSettings.FTCPAddress
}

func (p *sSettings) GetLogger() gopeer_logger.ILogger {
	return p.FSubSettings.FLogger
}
