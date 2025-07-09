package config

import (
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/app/config"
	"github.com/number571/hidden-lake/internal/utils/language"
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
