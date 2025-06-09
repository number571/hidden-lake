package http

import (
	"time"

	"github.com/number571/hidden-lake/pkg/adapters"
)

const (
	CDefaultReadTimeout   = 5 * time.Second
	CDefaultWriteTimeout  = 5 * time.Second
	CDefaultHandleTimeout = 5 * time.Second
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
	FAddress       string
	FReadTimeout   time.Duration
	FWriteTimeout  time.Duration
	FHandleTimeout time.Duration
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
	if p.FServeSettings.FReadTimeout == 0 {
		p.FServeSettings.FReadTimeout = CDefaultReadTimeout
	}
	if p.FServeSettings.FWriteTimeout == 0 {
		p.FServeSettings.FWriteTimeout = CDefaultWriteTimeout
	}
	if p.FServeSettings.FHandleTimeout == 0 {
		p.FServeSettings.FHandleTimeout = CDefaultHandleTimeout
	}
	return p
}

func (p *sSettings) GetAddress() string {
	return p.FServeSettings.FAddress
}

func (p *sSettings) GetAdapterSettings() adapters.ISettings {
	return p.FAdapterSettings
}

func (p *sSettings) GetReadTimeout() time.Duration {
	return p.FServeSettings.FReadTimeout
}

func (p *sSettings) GetWriteTimeout() time.Duration {
	return p.FServeSettings.FWriteTimeout
}

func (p *sSettings) GetHandleTimeout() time.Duration {
	return p.FServeSettings.FHandleTimeout
}
