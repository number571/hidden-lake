package config

import (
	"time"

	"github.com/number571/hidden-lake/internal/applications/remoter/pkg/app/config"
)

func GetConfigSettings(pCfg config.IConfig) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FExecTimeoutMS: uint64(sett.GetExecTimeout() / time.Millisecond), //nolint:gosec
			FPassword:      sett.GetPassword(),
		},
	}
}
