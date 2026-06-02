package config

import (
	"github.com/number571/hidden-lake/internal/adapters/meshtastic/pkg/app/config"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FMessageSizeBytes: sett.GetMessageSizeBytes(),
			FDatabaseEnabled:  sett.GetDatabaseEnabled(),
		},
	}
}
