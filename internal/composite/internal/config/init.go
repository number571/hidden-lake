package config

import (
	"os"

	hlf_settings "github.com/number571/hidden-lake/internal/applications/filesharer/pkg/settings"
	hlm_settings "github.com/number571/hidden-lake/internal/applications/messenger/pkg/settings"
	logger "github.com/number571/hidden-lake/internal/modules/logger/std"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		return LoadConfig(cfgPath)
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FServices: []string{
				hls_settings.CServiceFullName,
				hlm_settings.CServiceFullName,
				hlf_settings.CServiceFullName,
			},
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
