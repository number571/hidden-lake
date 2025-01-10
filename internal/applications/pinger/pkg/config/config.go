package config

import (
	"github.com/number571/hidden-lake/internal/applications/pinger/pkg/app/config"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FResponseMessage: sett.GetResponseMessage(),
		},
	}
}
