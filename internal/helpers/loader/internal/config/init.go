package config

import (
	"errors"
	"os"

	"github.com/number571/go-peer/pkg/encoding"
	hiddenlake "github.com/number571/hidden-lake"
	hll_settings "github.com/number571/hidden-lake/internal/helpers/loader/pkg/settings"
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
	network, ok := hiddenlake.GNetworks[pUseNetwork]
	if !ok {
		return nil, errors.Join(ErrRebuildConfig, ErrNetworkNotFound)
	}

	cfg.FSettings.FWorkSizeBits = network.FWorkSizeBits
	cfg.FSettings.FNetworkKey = pUseNetwork

	if err := os.WriteFile(cfg.fFilepath, encoding.SerializeYAML(cfg), 0o600); err != nil {
		return nil, errors.Join(ErrRebuildConfig, ErrWriteConfig, err)
	}

	rCfg, err := LoadConfig(cfg.fFilepath)
	if err != nil {
		return nil, errors.Join(ErrRebuildConfig, ErrLoadConfig, err)
	}

	return rCfg, nil
}

func initConfig() *SConfig {
	defaultNetwork := hiddenlake.GNetworks[hiddenlake.CDefaultNetwork]
	return &SConfig{
		FSettings: &SConfigSettings{
			FMessagesCapacity: defaultNetwork.FMessageSizeBytes,
			FWorkSizeBits:     defaultNetwork.FWorkSizeBits,
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
