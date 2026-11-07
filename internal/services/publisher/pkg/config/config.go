package config

import (
	"github.com/number571/hidden-lake/internal/services/publisher/pkg/app/config"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	_ = pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{},
	}
}
