package meshtastic

import (
	"time"

	"github.com/number571/hidden-lake/pkg/network/adapters"
)

var (
	_ ISettings = &sSettings{}
)

const (
	CDefaultWatchDuration = time.Second
	CDefaultReadTimeout   = 5 * time.Second
	CDefaultWriteTimeout  = 5 * time.Second
)

type SSettings sSettings
type sSettings struct {
	FAdapterSettings adapters.ISettings
	FServeSettings   *SServeSettings
}

type SServeSettings struct {
	FPath         string
	FAddress      string
	FDevPath      string
	FChannel      uint8
	FWatchPeriod  time.Duration
	FReadTimeout  time.Duration
	FWriteTimeout time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	if pSett == nil {
		pSett = &SSettings{
			FAdapterSettings: adapters.NewSettings(nil),
		}
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
	if p.FServeSettings.FPath == "" {
		p.FServeSettings.FPath = "."
	}
	if p.FServeSettings.FWatchPeriod == 0 {
		p.FServeSettings.FWatchPeriod = CDefaultWatchDuration
	}
	if p.FServeSettings.FReadTimeout == 0 {
		p.FServeSettings.FReadTimeout = CDefaultReadTimeout
	}
	if p.FServeSettings.FWriteTimeout == 0 {
		p.FServeSettings.FWriteTimeout = CDefaultWriteTimeout
	}
	return p
}

func (p *sSettings) GetAdapterSettings() adapters.ISettings {
	return p.FAdapterSettings
}

func (p *sSettings) GetPath() string {
	return p.FServeSettings.FPath
}

func (p *sSettings) GetAddress() string {
	return p.FServeSettings.FAddress
}

func (p *sSettings) GetDevPath() string {
	return p.FServeSettings.FDevPath
}

func (p *sSettings) GetChannel() uint8 {
	return p.FServeSettings.FChannel
}

func (p *sSettings) GetWatchPeriod() time.Duration {
	return p.FServeSettings.FWatchPeriod
}

func (p *sSettings) GetReadTimeout() time.Duration {
	return p.FServeSettings.FReadTimeout
}

func (p *sSettings) GetWriteTimeout() time.Duration {
	return p.FServeSettings.FWriteTimeout
}
