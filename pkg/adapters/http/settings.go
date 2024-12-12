package http

import (
	"github.com/number571/hidden-lake/pkg/adapters"
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FAddress         string
	FAdapterSettings adapters.ISettings
}

func NewSettings(pSett *SSettings) ISettings {
	if pSett == nil {
		pSett = &SSettings{}
		pSett.FAdapterSettings = adapters.NewSettings(&adapters.SSettings{})
	}
	return (&sSettings{
		FAddress:         pSett.FAddress,
		FAdapterSettings: pSett.FAdapterSettings,
	}).useDefault()
}

func (p *sSettings) useDefault() *sSettings {
	return p
}

func (p *sSettings) GetAddress() string {
	return p.FAddress
}

func (p *sSettings) GetAdapterSettings() adapters.ISettings {
	return p.FAdapterSettings
}
