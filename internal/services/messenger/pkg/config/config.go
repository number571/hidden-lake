package config

import (
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/app/config"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{},
	}
}
