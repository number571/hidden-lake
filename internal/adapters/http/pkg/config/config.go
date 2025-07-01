package config

import (
	"github.com/number571/hidden-lake/internal/adapters/http/pkg/app/config"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FMessageSizeBytes: sett.GetMessageSizeBytes(),
			FWorkSizeBits:     sett.GetMessageSizeBytes(),
			FNetworkKey:       sett.GetNetworkKey(),
			FDatabaseEnabled:  sett.GetDatabaseEnabled(),
			FReadTimeoutMS:    uint64(sett.GetReadTimeout().Milliseconds()),   // nolint: gosec
			FHandleTimeoutMS:  uint64(sett.GetHandleTimeout().Milliseconds()), // nolint: gosec
		},
	}
}
