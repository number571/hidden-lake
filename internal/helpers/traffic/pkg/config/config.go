package config

import (
	"github.com/number571/hidden-lake/internal/helpers/traffic/pkg/app/config"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FMessageSizeBytes: sett.GetMessageSizeBytes(),
			FWorkSizeBits:     sett.GetWorkSizeBits(),
			FMessagesCapacity: sett.GetMessagesCapacity(),
			FNetworkKey:       sett.GetNetworkKey(),
			FDatabaseEnabled:  sett.GetDatabaseEnabled(),
		},
	}
}
