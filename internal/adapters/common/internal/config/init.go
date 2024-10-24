package config

import (
	"os"

	hla_settings "github.com/number571/hidden-lake/internal/adapters/common/pkg/settings"
	hlt_settings "github.com/number571/hidden-lake/internal/helpers/traffic/pkg/settings"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		return LoadConfig(cfgPath)
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FSettings: &SConfigSettings{
				FWorkSizeBits: hls_settings.CDefaultWorkSizeBits,
				FNetworkKey:   hls_settings.CDefaultNetworkKey,
				FWaitTimeMS:   hla_settings.CDefaultWaitTimeMS,
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: hla_settings.CDefaultHTTPAddress,
			FConnection: &SConnection{
				FHLTHost: hlt_settings.CDefaultHTTPAddress,
				FSrvHost: hla_settings.CDefaultSrvAddress,
			},
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
