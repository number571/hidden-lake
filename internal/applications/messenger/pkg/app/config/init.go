package config

import (
	"os"

	logger "github.com/number571/hidden-lake/internal/utils/logger/std"

	hlm_settings "github.com/number571/hidden-lake/internal/applications/messenger/pkg/settings"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

func InitConfig(pCfgPath string, pInitCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(pCfgPath); !os.IsNotExist(err) {
		return LoadConfig(pCfgPath)
	}
	if pInitCfg == nil {
		pInitCfg = initConfig()
	}
	return BuildConfig(pCfgPath, pInitCfg)
}

func initConfig() *SConfig {
	return &SConfig{
		FSettings: &SConfigSettings{
			FMessagesCapacity: hlm_settings.CDefaultMessagesCapacity,
			FLanguage:         hlm_settings.CDefaultLanguage,
		},
		FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
		FAddress: &SAddress{
			FInternal: hlm_settings.CDefaultInternalAddress,
			FExternal: hlm_settings.CDefaultExternalAddress,
		},
		FConnection: hls_settings.CDefaultInternalAddress,
	}
}
