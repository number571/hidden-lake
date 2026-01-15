package config

import (
	"os"

	"github.com/number571/go-peer/pkg/crypto/random"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	hls_remoter_settings "github.com/number571/hidden-lake/internal/services/remoter/pkg/settings"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		return LoadConfig(cfgPath)
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FSettings: &SConfigSettings{
				FPassword: random.NewRandom().GetString(32),
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FInternal: hls_remoter_settings.CDefaultInternalAddress,
				FExternal: hls_remoter_settings.CDefaultExternalAddress,
			},
			FConnection: hlk_settings.CDefaultInternalAddress,
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
