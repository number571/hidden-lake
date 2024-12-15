package config

import (
	"errors"
	"net/url"
	"os"

	"github.com/number571/hidden-lake/build"
	hlm_settings "github.com/number571/hidden-lake/internal/applications/messenger/pkg/settings"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"

	hla_tcp_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
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

func initConfig() *SConfig {
	return &SConfig{
		FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
		FServices: []string{
			hla_tcp_settings.CServiceFullName,
			hls_settings.CServiceFullName,
			hlm_settings.CServiceFullName,
		},
	}
}

func rebuildConfig(pCfg IConfig, pUseNetwork string) (IConfig, error) {
	if pUseNetwork == "" {
		return pCfg, nil
	}

	cfg := pCfg.(*SConfig)
	network, ok := build.GNetworks[pUseNetwork]
	if !ok {
		return nil, errors.Join(ErrRebuildConfig, ErrNetworkNotFound)
	}

	mapAdapters := make(map[string]struct{}, 16)

	cfg.FServices = getServicesWithoutAdapters(pCfg)
	for _, c := range network.FConnections {
		u, err := url.Parse(c)
		if err != nil {
			return nil, errors.Join(ErrParseURL, err)
		}

		if _, ok := mapAdapters[u.Scheme]; ok {
			continue
		}
		mapAdapters[u.Scheme] = struct{}{}

		switch u.Scheme { // nolint: gocritic
		case hla_tcp_settings.CServiceAdapterScheme:
			cfg.FServices = append(cfg.FServices, hla_tcp_settings.CServiceFullName)
		}
	}

	return cfg, nil
}

func getServicesWithoutAdapters(cfg IConfig) []string {
	services := make([]string, 0, 16)
	for _, s := range cfg.GetServices() {
		switch s {
		case hla_tcp_settings.CServiceFullName:
			continue
		default:
			services = append(services, s)
		}
	}
	return services
}
