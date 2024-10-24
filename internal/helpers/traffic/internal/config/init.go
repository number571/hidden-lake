package config

import (
	"fmt"
	"os"

	logger "github.com/number571/hidden-lake/internal/utils/logger/std"

	hlt_settings "github.com/number571/hidden-lake/internal/helpers/traffic/pkg/settings"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		cfg, err := LoadConfig(cfgPath)
		if err != nil {
			return nil, fmt.Errorf("load config: %w", err)
		}
		return cfg, nil
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FSettings: &SConfigSettings{
				FMessageSizeBytes: hls_settings.CDefaultMessageSizeBytes,
				FWorkSizeBits:     hls_settings.CDefaultWorkSizeBits,
				FMessagesCapacity: hlt_settings.CDefaultMessagesCapacity,
				FDatabaseEnabled:  hlt_settings.CDefaultDatabaseEnabled,
				FNetworkKey:       hls_settings.CDefaultNetworkKey,
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FTCP:   hlt_settings.CDefaultTCPAddress,
				FHTTP:  hlt_settings.CDefaultHTTPAddress,
				FPPROF: "",
			},
			FConnections: []string{
				hlt_settings.CDefaultConnectionAddress,
			},
			FConsumers: []string{},
		}
	}
	cfg, err := BuildConfig(cfgPath, initCfg)
	if err != nil {
		return nil, fmt.Errorf("build config: %w", err)
	}
	return cfg, nil
}
