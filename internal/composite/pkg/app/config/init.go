package config

import (
	"os"

	hla_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
	hlm_settings "github.com/number571/hidden-lake/internal/applications/messenger/pkg/settings"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
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
				hla_settings.CServiceFullName,
			},
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
