package config

import (
	"github.com/number571/hidden-lake/internal/applications/filesharer/internal/config"
	"github.com/number571/hidden-lake/internal/modules/language"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FPageOffset: sett.GetPageOffset(),
			FRetryNum:   sett.GetRetryNum(),
			FLanguage:   language.FromILanguage(sett.GetLanguage()),
		},
	}
}
