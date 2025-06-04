package config

import (
	"errors"
	"os"

	hlm_settings "github.com/number571/hidden-lake/internal/applications/messenger/pkg/settings"
	hlp_settings "github.com/number571/hidden-lake/internal/applications/pinger/pkg/settings"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"

	hla_tcp_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
)

func InitConfig(cfgPath string, initCfg *SConfig, _ string) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		cfg, err := LoadConfig(cfgPath)
		if err != nil {
			return nil, errors.Join(ErrLoadConfig, err)
		}
		return cfg, nil
	}
	if initCfg == nil {
		initCfg = initConfig()
	}
	cfg, err := BuildConfig(cfgPath, initCfg)
	if err != nil {
		return nil, errors.Join(ErrBuildConfig, err)
	}
	return cfg, nil
}

func initConfig() *SConfig {
	return &SConfig{
		FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
		FServices: []string{
			hla_tcp_settings.CServiceFullName,
			hls_settings.CServiceFullName,
			hlp_settings.CServiceFullName,
			hlm_settings.CServiceFullName,
		},
	}
}
