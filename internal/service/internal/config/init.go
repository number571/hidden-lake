package config

import (
	"os"

	"github.com/number571/go-peer/pkg/utils"
	hlf_settings "github.com/number571/hidden-lake/internal/applications/filesharer/pkg/settings"
	hlm_settings "github.com/number571/hidden-lake/internal/applications/messenger/pkg/settings"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		cfg, err := LoadConfig(cfgPath)
		if err != nil {
			return nil, utils.MergeErrors(ErrLoadConfig, err)
		}
		return cfg, nil
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FSettings: &SConfigSettings{
				FMessageSizeBytes: hls_settings.CDefaultMessageSizeBytes,
				FWorkSizeBits:     hls_settings.CDefaultWorkSizeBits,
				FFetchTimeoutMS:   hls_settings.CDefaultFetchTimeoutMS,
				FQueuePeriodMS:    hls_settings.CDefaultQueuePeriodMS,
				FNetworkKey:       hls_settings.CDefaultNetworkKey,
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FTCP:  hls_settings.CDefaultTCPAddress,
				FHTTP: hls_settings.CDefaultHTTPAddress,
			},
			FServices: map[string]*SService{
				hlm_settings.CServiceFullName: {
					FHost: hlm_settings.CDefaultIncomingAddress,
				},
				hlf_settings.CServiceFullName: {
					FHost: hlf_settings.CDefaultIncomingAddress,
				},
			},
			FConnections: []string{},
			FFriends:     map[string]string{},
		}
	}
	cfg, err := BuildConfig(cfgPath, initCfg)
	if err != nil {
		return nil, utils.MergeErrors(ErrBuildConfig, err)
	}
	return cfg, nil
}
