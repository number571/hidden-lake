package config

import (
	"errors"
	"net/url"
	"os"
	"strings"

	"github.com/number571/hidden-lake/build"
	hlm_settings "github.com/number571/hidden-lake/internal/applications/messenger/pkg/settings"
	hlp_settings "github.com/number571/hidden-lake/internal/applications/pinger/pkg/settings"
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
			hlp_settings.CServiceFullName,
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

	mapAdapters := make(map[string]struct{}, len(network.FConnections))
	cfg.FServices = getServicesWithoutAdapters(pCfg)
	for _, c := range network.FConnections {
		u, err := url.Parse(c)
		if err != nil {
			return nil, errors.Join(ErrParseURL, err)
		}
		scheme := u.Scheme
		if _, ok := mapAdapters[scheme]; ok {
			continue
		}
		mapAdapters[scheme] = struct{}{}
		switch scheme { // nolint: gocritic
		case hla_tcp_settings.CServiceAdapterScheme:
			cfg.FServices = append(cfg.FServices, hla_tcp_settings.CServiceFullName)
		}
	}

	return cfg, nil
}

func getServicesWithoutAdapters(cfg IConfig) []string {
	services := cfg.GetServices()
	result := make([]string, 0, len(services))
	for _, s := range services {
		if strings.HasPrefix(s, "hidden-lake-adapter") {
			continue
		}
		result = append(result, s)
	}
	return result
}
