package config

import (
	"os"

	logger "github.com/number571/hidden-lake/internal/utils/logger/std"

	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	hls_messenger_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		return LoadConfig(cfgPath)
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FSettings: &SConfigSettings{
				FMessagesCapacity: hls_messenger_settings.CDefaultMessagesCapacity,
				FLanguage:         hls_messenger_settings.CDefaultLanguage,
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FInternal: hls_messenger_settings.CDefaultInternalAddress,
				FExternal: hls_messenger_settings.CDefaultExternalAddress,
			},
			FConnection: hlk_settings.CDefaultInternalAddress,
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
