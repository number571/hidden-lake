package config

import (
	"os"

	logger "github.com/number571/hidden-lake/internal/utils/logger/std"

	hlf_settings "github.com/number571/hidden-lake/internal/applications/filesharer/pkg/settings"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		return LoadConfig(cfgPath)
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FSettings: &SConfigSettings{
				FRetryNum:   hlf_settings.CDefaultRetryNum,
				FPageOffset: hlf_settings.CDefaultPageOffset,
				FLanguage:   hlf_settings.CDefaultLanguage,
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FInterface: hlf_settings.CDefaultInterfaceAddress,
				FIncoming:  hlf_settings.CDefaultIncomingAddress,
				FPPROF:     "",
			},
			FConnection: hls_settings.CDefaultHTTPAddress,
		}
	}
	return BuildConfig(cfgPath, initCfg)
}