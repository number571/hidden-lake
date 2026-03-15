package config

import (
	"github.com/number571/hidden-lake/internal/adapters/tcp/pkg/app/config"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	cfgSettings := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FMessageSizeBytes: cfgSettings.GetMessageSizeBytes(),
			FWorkSizeBits:     cfgSettings.GetMessageSizeBytes(),
			FNetworkKey:       cfgSettings.GetNetworkKey(),
			FDatabaseEnabled:  cfgSettings.GetDatabaseEnabled(),
			FConnKeepPeriodMS: uint64(cfgSettings.GetConnKeepPeriod().Milliseconds()), // nolint: gosec
			FSendTimeoutMS:    uint64(cfgSettings.GetSendTimeout().Milliseconds()),    // nolint: gosec
			FRecvTimeoutMS:    uint64(cfgSettings.GetRecvTimeout().Milliseconds()),    // nolint: gosec
			FDialTimeoutMS:    uint64(cfgSettings.GetDialTimeout().Milliseconds()),    // nolint: gosec
			FWaitTimeoutMS:    uint64(cfgSettings.GetWaitTimeout().Milliseconds()),    // nolint: gosec
		},
	}
}
