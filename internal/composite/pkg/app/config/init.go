package config

import (
	"errors"
	"os"

	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	hls_messenger_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	hls_pinger_settings "github.com/number571/hidden-lake/internal/services/pinger/pkg/settings"
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
		FApplications: []string{
			hlk_settings.CAppShortName,
			hla_tcp_settings.CAppShortName,
			hls_pinger_settings.CAppShortName,
			hls_messenger_settings.CAppShortName,
			hls_filesharer_settings.CAppShortName,
		},
	}
}
