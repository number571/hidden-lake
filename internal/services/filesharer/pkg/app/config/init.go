package config

import (
	"os"

	logger "github.com/number571/hidden-lake/internal/utils/logger/std"

	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		return LoadConfig(cfgPath)
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FSettings: &SConfigSettings{
				FRetryNum:   hls_filesharer_settings.CDefaultRetryNum,
				FPageOffset: hls_filesharer_settings.CDefaultPageOffset,
				FLanguage:   hls_filesharer_settings.CDefaultLanguage,
			},
			FLogging:    []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress:    hls_filesharer_settings.CDefaultAddress,
			FConnection: hlk_settings.CDefaultInternalAddress,
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
