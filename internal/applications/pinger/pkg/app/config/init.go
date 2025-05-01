package config

import (
	"os"

	hlp_settings "github.com/number571/hidden-lake/internal/applications/pinger/pkg/settings"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
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
		FSettings: &SConfigSettings{},
		FLogging:  []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
		FAddress: &SAddress{
			FExternal: hlp_settings.CDefaultExternalAddress,
		},
	}
}
