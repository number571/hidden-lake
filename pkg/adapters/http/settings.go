package http

import (
	"github.com/number571/hidden-lake/pkg/adapters"
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FAdapterSettings adapters.ISettings
	FSrvSettings     *SSrvSettings
}

type SSrvSettings struct {
	FAddress string
}

func NewSettings(pSett *SSettings) ISettings {
	if pSett == nil {
		pSett = &SSettings{FAdapterSettings: adapters.NewSettings(nil)}
	}
	return (&sSettings{
		FAdapterSettings: pSett.FAdapterSettings,
		FSrvSettings:     pSett.FSrvSettings,
	}).initDefault()
}

func (p *sSettings) initDefault() *sSettings {
	if p.FSrvSettings == nil {
		p.FSrvSettings = &SSrvSettings{}
	}
	return p
}

func (p *sSettings) GetAddress() string {
	return p.FSrvSettings.FAddress
}

func (p *sSettings) GetAdapterSettings() adapters.ISettings {
	return p.FAdapterSettings
}
