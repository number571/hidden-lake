package config

import (
	"errors"
	"net/url"
	"os"

	"github.com/number571/hidden-lake/build"
	hla_http_settings "github.com/number571/hidden-lake/internal/adapters/http/pkg/settings"
	hla_https_settings "github.com/number571/hidden-lake/internal/adapters/https/pkg/settings"
	hla_tcp_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	hls_messenger_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	hls_pinger_settings "github.com/number571/hidden-lake/internal/services/pinger/pkg/settings"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func InitConfig(cfgPath string, initCfg *SConfig, networkKey string) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		cfg, err := LoadConfig(cfgPath)
		if err != nil {
			return nil, errors.Join(ErrLoadConfig, err)
		}
		return cfg, nil
	}
	if initCfg == nil {
		cfg, err := initConfig(networkKey)
		if err != nil {
			return nil, errors.Join(ErrInitDefaultConfig, err)
		}
		initCfg = cfg
	}
	cfg, err := BuildConfig(cfgPath, initCfg)
	if err != nil {
		return nil, errors.Join(ErrBuildConfig, err)
	}
	return cfg, nil
}

func initConfig(networkKey string) (*SConfig, error) {
	defaultConfig := &SConfig{
		FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
		FApplications: []string{
			hlk_settings.CAppShortName,
			hls_pinger_settings.CAppShortName,
			hls_messenger_settings.CAppShortName,
			hls_filesharer_settings.CAppShortName,
		},
	}

	network, ok := build.GetNetwork(networkKey)
	if !ok {
		return defaultConfig, nil
	}

	mapUsed := make(map[string]struct{}, 16)
	for _, c := range network.FConnections {
		u, err := url.Parse(c)
		if err != nil {
			return nil, ErrParseConnection
		}

		scheme := u.Scheme
		if _, ok := mapUsed[scheme]; ok {
			continue
		}
		mapUsed[scheme] = struct{}{}

		adapterName := ""
		switch scheme {
		case hla_tcp_settings.CAppAdapterName:
			adapterName = hla_tcp_settings.CAppShortName
		case hla_http_settings.CAppAdapterName:
			adapterName = hla_http_settings.CAppShortName
		case hla_https_settings.CAppAdapterName:
			adapterName = hla_https_settings.CAppShortName
		}

		if adapterName == "" {
			return nil, ErrAdapterNotFound
		}
		defaultConfig.FApplications = append(defaultConfig.FApplications, adapterName)
	}

	return defaultConfig, nil
}
