package config

import (
	"os"

	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	hls_pinger_settings "github.com/number571/hidden-lake/internal/services/pinger/pkg/settings"
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
				FInternal: hls_pinger_settings.CDefaultInternalAddress,
				FExternal: hls_pinger_settings.CDefaultExternalAddress,
			},
			FConnection: hlk_settings.CDefaultInternalAddress,
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
