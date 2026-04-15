package https

import (
	"time"

	"github.com/number571/hidden-lake/pkg/network/adapters"
)

const (
	CDefaultChannelSize   = 64
	CDefaultConnNumLimit  = 256
	CDefaultReadTimeout   = 5 * time.Second
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
	FAuthMapper    map[string]string
	FChannelSize   uint64
	FConnNumLimit  uint64
	FReadTimeout   time.Duration
	FHandleTimeout time.Duration
}

// TODO: authMapper
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
	if p.FServeSettings.FChannelSize == 0 {
		p.FServeSettings.FChannelSize = CDefaultChannelSize
	}
	if p.FServeSettings.FConnNumLimit == 0 {
		p.FServeSettings.FConnNumLimit = CDefaultConnNumLimit
	}
	if p.FServeSettings.FReadTimeout == 0 {
		p.FServeSettings.FReadTimeout = CDefaultReadTimeout
	}
	if p.FServeSettings.FHandleTimeout == 0 {
		p.FServeSettings.FHandleTimeout = CDefaultHandleTimeout
	}
	return p
}

func (p *sSettings) GetAddress() string {
	return p.FServeSettings.FAddress
}

func (p *sSettings) GetAuthMapper() map[string]string {
	return p.FServeSettings.FAuthMapper
}

func (p *sSettings) GetChannelSize() uint64 {
	return p.FServeSettings.FChannelSize
}

func (p *sSettings) GetConnNumLimit() uint64 {
	return p.FServeSettings.FConnNumLimit
}

func (p *sSettings) GetAdapterSettings() adapters.ISettings {
	return p.FAdapterSettings
}

func (p *sSettings) GetReadTimeout() time.Duration {
	return p.FServeSettings.FReadTimeout
}

func (p *sSettings) GetHandleTimeout() time.Duration {
	return p.FServeSettings.FHandleTimeout
}
