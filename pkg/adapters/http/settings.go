package http

import (
	"time"

	"github.com/number571/hidden-lake/pkg/adapters"
)

const (
	CDefaultRecvTimeout = 5 * time.Second
	CDefaultSendTimeout = 5 * time.Second
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FAdapterSettings adapters.ISettings
	FServeSettings   *SServeSettings
}

type SServeSettings struct {
	FAddress     string
	FRecvTimeout time.Duration
	FSendTimeout time.Duration
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
	if p.FServeSettings.FRecvTimeout == 0 {
		p.FServeSettings.FRecvTimeout = CDefaultRecvTimeout
	}
	if p.FServeSettings.FSendTimeout == 0 {
		p.FServeSettings.FSendTimeout = CDefaultSendTimeout
	}
	return p
}

func (p *sSettings) GetAddress() string {
	return p.FServeSettings.FAddress
}

func (p *sSettings) GetAdapterSettings() adapters.ISettings {
	return p.FAdapterSettings
}

func (p *sSettings) GetRecvTimeout() time.Duration {
	return p.FServeSettings.FRecvTimeout
}

func (p *sSettings) GetSendTimeout() time.Duration {
	return p.FServeSettings.FSendTimeout
}
