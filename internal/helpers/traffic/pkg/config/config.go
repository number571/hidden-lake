package config

import (
	"github.com/number571/hidden-lake/internal/helpers/traffic/internal/config"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FMessageSizeBytes:     sett.GetMessageSizeBytes(),
			FKeySizeBits:          sett.GetKeySizeBits(),
			FWorkSizeBits:         sett.GetWorkSizeBits(),
			FMessagesCapacity:     sett.GetMessagesCapacity(),
			FRandMessageSizeBytes: sett.GetRandMessageSizeBytes(),
			FNetworkKey:           sett.GetNetworkKey(),
			FDatabaseEnabled:      sett.GetDatabaseEnabled(),
		},
	}
}