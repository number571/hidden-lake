package tcp

import (
	"time"

	"github.com/number571/hidden-lake/pkg/adapters"
)

var (
	_ ISettings = &sSettings{}
)

const (
	CDefaultConnLimit      = 256
	CDefaultConnKeepPeriod = 10 * time.Second
	CDefaultSendTimeout    = 5 * time.Second
	CDefaultRecvTimeout    = 5 * time.Second
	CDefaultDialTimeout    = 5 * time.Second
	CDefaultWaitTimeout    = time.Hour
)

type SSettings sSettings
type sSettings struct {
	FAdapterSettings adapters.ISettings
	FSrvSettings     *SSrvSettings
}

type SSrvSettings struct {
	FAddress        string
	FConnLimit      uint64
	FConnKeepPeriod time.Duration
	FSendTimeout    time.Duration
	FRecvTimeout    time.Duration
	FDialTimeout    time.Duration
	FWaitTimeout    time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	if pSett == nil {
		pSett = &SSettings{}
	}
	return (&sSettings{
		FAdapterSettings: pSett.FAdapterSettings,
		FSrvSettings:     pSett.FSrvSettings,
	}).useDefault()
}

func (p *sSettings) useDefault() *sSettings {
	if p.FSrvSettings == nil {
		p.FSrvSettings = &SSrvSettings{}
	}
	if p.FSrvSettings.FConnLimit == 0 {
		p.FSrvSettings.FConnLimit = CDefaultConnLimit
	}
	if p.FSrvSettings.FConnKeepPeriod == 0 {
		p.FSrvSettings.FConnKeepPeriod = CDefaultConnKeepPeriod
	}
	if p.FSrvSettings.FSendTimeout == 0 {
		p.FSrvSettings.FSendTimeout = CDefaultSendTimeout
	}
	if p.FSrvSettings.FRecvTimeout == 0 {
		p.FSrvSettings.FRecvTimeout = CDefaultRecvTimeout
	}
	if p.FSrvSettings.FDialTimeout == 0 {
		p.FSrvSettings.FDialTimeout = CDefaultDialTimeout
	}
	if p.FSrvSettings.FWaitTimeout == 0 {
		p.FSrvSettings.FWaitTimeout = CDefaultWaitTimeout
	}
	return p
}

func (p *sSettings) GetAddress() string {
	return p.FSrvSettings.FAddress
}

func (p *sSettings) GetConnLimit() uint64 {
	return p.FSrvSettings.FConnLimit
}

func (p *sSettings) GetConnKeepPeriod() time.Duration {
	return p.FSrvSettings.FConnKeepPeriod
}

func (p *sSettings) GetSendTimeout() time.Duration {
	return p.FSrvSettings.FSendTimeout
}

func (p *sSettings) GetRecvTimeout() time.Duration {
	return p.FSrvSettings.FRecvTimeout
}

func (p *sSettings) GetDialTimeout() time.Duration {
	return p.FSrvSettings.FDialTimeout
}

func (p *sSettings) GetWaitTimeout() time.Duration {
	return p.FSrvSettings.FWaitTimeout
}

func (p *sSettings) GetAdapterSettings() adapters.ISettings {
	return p.FAdapterSettings
}
