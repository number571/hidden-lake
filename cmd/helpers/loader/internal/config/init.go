package config

import (
	"os"

	hll_settings "github.com/number571/hidden-lake/cmd/helpers/loader/pkg/settings"
	hls_settings "github.com/number571/hidden-lake/cmd/service/pkg/settings"
	logger "github.com/number571/hidden-lake/internal/logger/std"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		return LoadConfig(cfgPath)
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FSettings: &SConfigSettings{
				FMessagesCapacity: hll_settings.CDefaultMessagesCapacity,
				FWorkSizeBits:     hls_settings.CDefaultWorkSizeBits,
				FNetworkKey:       hls_settings.CDefaultNetworkKey,
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FHTTP:  hll_settings.CDefaultHTTPAddress,
				FPPROF: "",
			},
			FProducers: []string{},
			FConsumers: []string{},
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
