package tcp

import (
	"time"

	"github.com/number571/hidden-lake/pkg/adapters"
)

var (
	_ ISettings = &sSettings{}
)

const (
	CDefaultConnNumLimit   = 256
	CDefaultConnKeepPeriod = 10 * time.Second
	CDefaultSendTimeout    = 5 * time.Second
	CDefaultRecvTimeout    = 5 * time.Second
	CDefaultDialTimeout    = 5 * time.Second
	CDefaultWaitTimeout    = time.Hour
)

type SSettings sSettings
type sSettings struct {
	FAdapterSettings adapters.ISettings
	FServeSettings   *SServeSettings
}

type SServeSettings struct {
	FAddress        string
	FConnNumLimit   uint64
	FConnKeepPeriod time.Duration
	FSendTimeout    time.Duration
	FRecvTimeout    time.Duration
	FDialTimeout    time.Duration
	FWaitTimeout    time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	if pSett == nil {
		pSett = &SSettings{FAdapterSettings: adapters.NewSettings(nil)}
	}
	return (&sSettings{
		FAdapterSettings: pSett.FAdapterSettings,
		FServeSettings:   pSett.FServeSettings,
	}).initDefault()
}

func (p *sSettings) initDefault() *sSettings {
	if p.FServeSettings == nil {
		p.FServeSettings = &SServeSettings{}
	}
	if p.FServeSettings.FConnNumLimit == 0 {
		p.FServeSettings.FConnNumLimit = CDefaultConnNumLimit
	}
	if p.FServeSettings.FConnKeepPeriod == 0 {
		p.FServeSettings.FConnKeepPeriod = CDefaultConnKeepPeriod
	}
	if p.FServeSettings.FSendTimeout == 0 {
		p.FServeSettings.FSendTimeout = CDefaultSendTimeout
	}
	if p.FServeSettings.FRecvTimeout == 0 {
		p.FServeSettings.FRecvTimeout = CDefaultRecvTimeout
	}
	if p.FServeSettings.FDialTimeout == 0 {
		p.FServeSettings.FDialTimeout = CDefaultDialTimeout
	}
	if p.FServeSettings.FWaitTimeout == 0 {
		p.FServeSettings.FWaitTimeout = CDefaultWaitTimeout
	}
	return p
}

func (p *sSettings) GetAddress() string {
	return p.FServeSettings.FAddress
}

func (p *sSettings) GetConnNumLimit() uint64 {
	return p.FServeSettings.FConnNumLimit
}

func (p *sSettings) GetConnKeepPeriod() time.Duration {
	return p.FServeSettings.FConnKeepPeriod
}

func (p *sSettings) GetSendTimeout() time.Duration {
	return p.FServeSettings.FSendTimeout
}

func (p *sSettings) GetRecvTimeout() time.Duration {
	return p.FServeSettings.FRecvTimeout
}

func (p *sSettings) GetDialTimeout() time.Duration {
	return p.FServeSettings.FDialTimeout
}

func (p *sSettings) GetWaitTimeout() time.Duration {
	return p.FServeSettings.FWaitTimeout
}

func (p *sSettings) GetAdapterSettings() adapters.ISettings {
	return p.FAdapterSettings
}
