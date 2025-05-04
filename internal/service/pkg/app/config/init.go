package config

import (
	"errors"
	"net/url"
	"os"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/hidden-lake/build"
	hla_tcp_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
	hlm_settings "github.com/number571/hidden-lake/internal/applications/messenger/pkg/settings"
	hlp_settings "github.com/number571/hidden-lake/internal/applications/pinger/pkg/settings"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func InitConfig(pCfgPath string, pInitCfg *SConfig, pUseNetwork string) (IConfig, error) {
	if _, err := os.Stat(pCfgPath); !os.IsNotExist(err) {
		cfg, err := LoadConfig(pCfgPath)
		if err != nil {
			return nil, errors.Join(ErrLoadConfig, err)
		}
		return rebuildConfig(cfg, pUseNetwork)
	}
	if pInitCfg == nil {
		pInitCfg = initConfig()
	}
	cfg, err := BuildConfig(pCfgPath, pInitCfg)
	if err != nil {
		return nil, errors.Join(ErrBuildConfig, err)
	}
	return rebuildConfig(cfg, pUseNetwork)
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

	cfg.FSettings.FFetchTimeoutMS = network.FFetchTimeoutMS
	cfg.FSettings.FMessageSizeBytes = network.FMessageSizeBytes
	cfg.FSettings.FQueuePeriodMS = network.FQueuePeriodMS
	cfg.FSettings.FWorkSizeBits = network.FWorkSizeBits
	cfg.FSettings.FNetworkKey = pUseNetwork

	cfg.FEndpoints = make([]string, 0, len(network.FConnections))
	for _, c := range network.FConnections {
		u, err := url.Parse(c)
		if err != nil {
			return nil, errors.Join(ErrParseURL, err)
		}
		if u.Scheme != hls_settings.CServiceAdapterScheme {
			continue
		}
		cfg.FEndpoints = append(cfg.FEndpoints, u.Host)
	}

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
	defaultNetwork, _ := build.GetNetwork(build.CDefaultNetwork)
	return &SConfig{
		FSettings: &SConfigSettings{
			FMessageSizeBytes: defaultNetwork.FMessageSizeBytes,
			FWorkSizeBits:     defaultNetwork.FWorkSizeBits,
			FFetchTimeoutMS:   defaultNetwork.FFetchTimeoutMS,
			FQueuePeriodMS:    defaultNetwork.FQueuePeriodMS,
		},
		FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
		FAddress: &SAddress{
			FExternal: hls_settings.CDefaultExternalAddress,
			FInternal: hls_settings.CDefaultInternalAddress,
		},
		FServices: map[string]string{
			hlp_settings.CServiceFullName: hlp_settings.CDefaultExternalAddress,
			hlm_settings.CServiceFullName: hlm_settings.CDefaultExternalAddress,
		},
		FEndpoints: []string{
			hla_tcp_settings.CDefaultInternalAddress,
		},
		FFriends: map[string]string{},
	}
}
