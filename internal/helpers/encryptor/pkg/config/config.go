package config

import (
	"github.com/number571/hidden-lake/internal/helpers/encryptor/internal/config"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FMessageSizeBytes: sett.GetMessageSizeBytes(),
			FWorkSizeBits:     sett.GetWorkSizeBits(),
			FNetworkKey:       sett.GetNetworkKey(),
		},
	}
}
