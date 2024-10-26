package config

import (
	"os"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/utils"
	hiddenlake "github.com/number571/hidden-lake"
	hlf_settings "github.com/number571/hidden-lake/internal/applications/filesharer/pkg/settings"
	hlm_settings "github.com/number571/hidden-lake/internal/applications/messenger/pkg/settings"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
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

func rebuildConfig(pCfg IConfig, pUseNetwork string) (IConfig, error) {
	if pUseNetwork == "" {
		return pCfg, nil
	}

	cfg := pCfg.(*SConfig)
	network, ok := hiddenlake.GNetworks[pUseNetwork]
	if !ok {
		return nil, utils.MergeErrors(ErrRebuildConfig, ErrNetworkNotFound)
	}

	cfg.FSettings.FFetchTimeoutMS = network.FFetchTimeoutMS
	cfg.FSettings.FMessageSizeBytes = network.FMessageSizeBytes
	cfg.FSettings.FQueuePeriodMS = network.FQueuePeriodMS
	cfg.FSettings.FWorkSizeBits = network.FWorkSizeBits
	cfg.FSettings.FNetworkKey = pUseNetwork
	cfg.FConnections = network.FConnections

	if err := os.WriteFile(cfg.fFilepath, encoding.SerializeYAML(cfg), 0o600); err != nil {
		return nil, utils.MergeErrors(ErrRebuildConfig, ErrWriteConfig, err)
	}

	rCfg, err := LoadConfig(cfg.fFilepath)
	if err != nil {
		return nil, utils.MergeErrors(ErrRebuildConfig, ErrLoadConfig, err)
	}

	return rCfg, nil
}

func initConfig() *SConfig {
	return &SConfig{
		FSettings: &SConfigSettings{
			FMessageSizeBytes: hls_settings.CDefaultMessageSizeBytes,
			FWorkSizeBits:     hls_settings.CDefaultWorkSizeBits,
			FFetchTimeoutMS:   hls_settings.CDefaultFetchTimeoutMS,
			FQueuePeriodMS:    hls_settings.CDefaultQueuePeriodMS,
			FNetworkKey:       hls_settings.CDefaultNetworkKey,
		},
		FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
		FAddress: &SAddress{
			FTCP:  hls_settings.CDefaultTCPAddress,
			FHTTP: hls_settings.CDefaultHTTPAddress,
		},
		FServices: map[string]string{
			hlm_settings.CServiceFullName: hlm_settings.CDefaultIncomingAddress,
			hlf_settings.CServiceFullName: hlf_settings.CDefaultIncomingAddress,
		},
		FConnections: []string{},
		FFriends:     map[string]string{},
	}
}
