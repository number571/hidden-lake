package config

import (
	"github.com/number571/hidden-lake/internal/applications/notifier/pkg/app/config"
	"github.com/number571/hidden-lake/internal/utils/language"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FMessagesCapacity: sett.GetMessagesCapacity(),
			FWorkSizeBits:     sett.GetWorkSizeBits(),
			FPowParallel:      sett.GetPowParallel(),
			FNetworkKey:       sett.GetNetworkKey(),
			FLanguage:         language.FromILanguage(sett.GetLanguage()),
		},
	}
}
