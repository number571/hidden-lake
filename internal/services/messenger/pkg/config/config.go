package config

import (
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/app/config"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FMessagesCapacity: sett.GetMessagesCapacity(),
		},
	}
}
