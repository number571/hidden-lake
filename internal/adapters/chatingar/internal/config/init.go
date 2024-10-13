package config

import (
	"os"

	hla_settings "github.com/number571/hidden-lake/internal/adapters/chatingar/pkg/settings"
	hlt_settings "github.com/number571/hidden-lake/internal/helpers/traffic/pkg/settings"
	logger "github.com/number571/hidden-lake/internal/modules/logger/std"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
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
				FPostID:  hla_settings.CDefaultPostID,
			},
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
