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
			FWatchPeriodMS:    uint64(sett.GetWatchPeriod().Milliseconds()),  // nolint:gosec
			FReadTimeoutMS:    uint64(sett.GetReadTimeout().Milliseconds()),  // nolint:gosec
			FWriteTimeoutMS:   uint64(sett.GetWriteTimeout().Milliseconds()), // nolint:gosec
			FMaxDelayTimeMS:   uint64(sett.GetMaxDelayTime().Milliseconds()), // nolint:gosec
		},
	}
}
