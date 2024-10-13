package config

import (
	"github.com/number571/hidden-lake/internal/applications/remoter/internal/config"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FExecTimeoutMS: sett.GetExecTimeoutMS(),
			FPassword:      sett.GetPassword(),
		},
	}
}
