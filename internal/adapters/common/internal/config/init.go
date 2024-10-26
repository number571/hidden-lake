package config

import (
	"os"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/utils"
	hiddenlake "github.com/number571/hidden-lake"
	hla_settings "github.com/number571/hidden-lake/internal/adapters/common/pkg/settings"
	hlt_settings "github.com/number571/hidden-lake/internal/helpers/traffic/pkg/settings"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func InitConfig(cfgPath string, initCfg *SConfig, useNetwork string) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		cfg, err := LoadConfig(cfgPath)
		if err != nil {
			return nil, utils.MergeErrors(ErrLoadConfig, err)
		}
		return rebuildConfig(cfg, useNetwork)
	}
	if initCfg == nil {
		initCfg = initConfig()
	}
	cfg, err := BuildConfig(cfgPath, initCfg)
	if err != nil {
		return nil, utils.MergeErrors(ErrBuildConfig, err)
	}
	return rebuildConfig(cfg, useNetwork)
}

func initConfig() *SConfig {
	defaultNetwork := hiddenlake.GNetworks[hiddenlake.CDefaultNetwork]
	return &SConfig{
		FSettings: &SConfigSettings{
			FWorkSizeBits: defaultNetwork.FWorkSizeBits,
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

func rebuildConfig(pCfg IConfig, pUseNetwork string) (IConfig, error) {
	if pUseNetwork == "" {
		return pCfg, nil
	}

	cfg := pCfg.(*SConfig)
	network, ok := hiddenlake.GNetworks[pUseNetwork]
	if !ok {
		return nil, utils.MergeErrors(ErrRebuildConfig, ErrNetworkNotFound)
	}

	cfg.FSettings.FWorkSizeBits = network.FWorkSizeBits
	cfg.FSettings.FNetworkKey = pUseNetwork

	if err := os.WriteFile(cfg.fFilepath, encoding.SerializeYAML(cfg), 0o600); err != nil {
		return nil, utils.MergeErrors(ErrRebuildConfig, ErrWriteConfig, err)
	}

	rCfg, err := LoadConfig(cfg.fFilepath)
	if err != nil {
		return nil, utils.MergeErrors(ErrRebuildConfig, ErrLoadConfig, err)
	}

	return rCfg, nil
}
