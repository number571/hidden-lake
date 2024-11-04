package config

import (
	"os"

	"github.com/number571/go-peer/pkg/crypto/random"
	hlr_settings "github.com/number571/hidden-lake/internal/applications/remoter/pkg/settings"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func InitConfig(cfgPath string, initCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		return LoadConfig(cfgPath)
	}
	if initCfg == nil {
		initCfg = &SConfig{
			FSettings: &SConfigSettings{
				FExecTimeoutMS: hlr_settings.CDefaultExecTimeout,
				FPassword:      random.NewRandom().GetString(32),
			},
			FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
			FAddress: &SAddress{
				FIncoming: hlr_settings.CDefaultIncomingAddress,
				FPPROF:    "",
			},
		}
	}
	return BuildConfig(cfgPath, initCfg)
}
