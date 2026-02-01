package config

import (
	"errors"
	"os"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/hidden-lake/build"
	hla_tcp_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	hls_messenger_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	hls_pinger_settings "github.com/number571/hidden-lake/internal/services/pinger/pkg/settings"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func InitConfig(cfgPath string, initCfg *SConfig, useNetwork string) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		cfg, err := LoadConfig(cfgPath)
		if err != nil {
			return nil, errors.Join(ErrLoadConfig, err)
		}
		return rebuildConfig(cfg, useNetwork)
	}
	if initCfg == nil {
		initCfg = initConfig()
	}
	cfg, err := BuildConfig(cfgPath, initCfg)
	if err != nil {
		return nil, errors.Join(ErrBuildConfig, err)
	}
	return rebuildConfig(cfg, useNetwork)
}

func rebuildConfig(pCfg IConfig, pUseNetwork string) (IConfig, error) {
	if pUseNetwork == "" {
		return pCfg, nil
	}

	cfg := pCfg.(*SConfig)
	network, ok := build.GetNetwork(pUseNetwork)
	if !ok {
		return nil, errors.Join(ErrRebuildConfig, ErrNetworkNotFound)
	}

	cfg.FSettings.FMessageSizeBytes = network.FMessageSizeBytes
	cfg.FSettings.FWorkSizeBits = network.FWorkSizeBits
	cfg.FSettings.FNetworkKey = pUseNetwork

	if err := os.WriteFile(cfg.fFilepath, encoding.SerializeYAML(cfg), 0600); err != nil {
		return nil, errors.Join(ErrRebuildConfig, ErrWriteConfig, err)
	}

	rCfg, err := LoadConfig(cfg.fFilepath)
	if err != nil {
		return nil, errors.Join(ErrRebuildConfig, ErrLoadConfig, err)
	}

	return rCfg, nil
}

func initConfig() *SConfig {
	return &SConfig{
		FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
		FAddress: &SAddress{
			FExternal: hlk_settings.CDefaultExternalAddress,
			FInternal: hlk_settings.CDefaultInternalAddress,
		},
		FServices: map[string]string{
			hls_pinger_settings.CAppShortName:    hls_pinger_settings.CDefaultExternalAddress,
			hls_messenger_settings.CAppShortName: hls_messenger_settings.CDefaultExternalAddress,
		},
		FEndpoints: []string{
			hla_tcp_settings.CDefaultInternalAddress,
		},
	}
}
