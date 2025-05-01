package config

import (
	"os"

	logger "github.com/number571/hidden-lake/internal/utils/logger/std"

	hlf_settings "github.com/number571/hidden-lake/internal/applications/filesharer/pkg/settings"
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
			FRetryNum:   hlf_settings.CDefaultRetryNum,
			FPageOffset: hlf_settings.CDefaultPageOffset,
			FLanguage:   hlf_settings.CDefaultLanguage,
		},
		FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
		FAddress: &SAddress{
			FInternal: hlf_settings.CDefaultInternalAddress,
			FExternal: hlf_settings.CDefaultExternalAddress,
		},
		FConnection: hls_settings.CDefaultInternalAddress,
	}
}
