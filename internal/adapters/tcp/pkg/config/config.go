package config

import (
	"github.com/number571/hidden-lake/internal/adapters/tcp/pkg/app/config"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FMessageSizeBytes: sett.GetMessageSizeBytes(),
			FWorkSizeBits:     sett.GetMessageSizeBytes(),
			FNetworkKey:       sett.GetNetworkKey(),
			FDatabaseEnabled:  sett.GetDatabaseEnabled(),
			FConnNumLimit:     sett.GetConnNumLimit(),
			FConnKeepPeriodMS: uint64(sett.GetConnKeepPeriod().Milliseconds()), // nolint: gosec
			FSendTimeoutMS:    uint64(sett.GetSendTimeout().Milliseconds()),    // nolint: gosec
			FRecvTimeoutMS:    uint64(sett.GetRecvTimeout().Milliseconds()),    // nolint: gosec
			FDialTimeoutMS:    uint64(sett.GetDialTimeout().Milliseconds()),    // nolint: gosec
			FWaitTimeoutMS:    uint64(sett.GetWaitTimeout().Milliseconds()),    // nolint: gosec
		},
	}
}
