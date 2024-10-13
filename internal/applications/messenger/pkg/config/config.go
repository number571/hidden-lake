package config

import (
	"github.com/number571/hidden-lake/internal/applications/messenger/internal/config"
	"github.com/number571/hidden-lake/internal/utils/language"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FMessagesCapacity: sett.GetMessagesCapacity(),
			FLanguage:         language.FromILanguage(sett.GetLanguage()),
		},
	}
}
