package config

import (
	"os"

	"github.com/number571/go-peer/pkg/crypto/random"
	hlr_settings "github.com/number571/hidden-lake/internal/applications/remoter/pkg/settings"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func InitConfig(pCfgPath string, pInitCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(pCfgPath); !os.IsNotExist(err) {
		return LoadConfig(pCfgPath)
	}
	if pInitCfg == nil {
		pInitCfg = initConfig()
	}
	return BuildConfig(pCfgPath, pInitCfg)
}

func initConfig() *SConfig {
	return &SConfig{
		FSettings: &SConfigSettings{
			FExecTimeoutMS: hlr_settings.CDefaultExecTimeout,
			FPassword:      random.NewRandom().GetString(32),
		},
		FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
		FAddress: &SAddress{
			FExternal: hlr_settings.CDefaultExternalAddress,
		},
	}
}
