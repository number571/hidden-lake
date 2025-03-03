package config

import (
	"os"

	hlp_settings "github.com/number571/hidden-lake/internal/applications/pinger/pkg/settings"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		return LoadConfig(cfgPath)
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FSettings: &SConfigSettings{},
			FLogging:  []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FExternal: hlp_settings.CDefaultExternalAddress,
			},
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
