package config

import (
	"os"

	hle_settings "github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/settings"
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
				FMessageSizeBytes: hls_settings.CDefaultMessageSizeBytes,
				FWorkSizeBits:     hls_settings.CDefaultWorkSizeBits,
				FNetworkKey:       hls_settings.CDefaultNetworkKey,
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FHTTP:  hle_settings.CDefaultHTTPAddress,
				FPPROF: "",
			},
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
