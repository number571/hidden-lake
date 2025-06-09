package config

import (
	"github.com/number571/hidden-lake/internal/applications/remoter/pkg/app/config"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FPassword: sett.GetPassword(),
		},
	}
}
