package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/number571/go-peer/pkg/encoding"
	hiddenlake "github.com/number571/hidden-lake"
	hlt_settings "github.com/number571/hidden-lake/internal/helpers/traffic/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/conn"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func InitConfig(cfgPath string, initCfg *SConfig, useNetwork string) (IConfig, error) {
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		cfg, err := LoadConfig(cfgPath)
		if err != nil {
			return nil, fmt.Errorf("load config: %w", err)
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

	cfg.FSettings.FMessageSizeBytes = network.FMessageSizeBytes
	cfg.FSettings.FWorkSizeBits = network.FWorkSizeBits
	cfg.FSettings.FNetworkKey = pUseNetwork

	cfg.FConnections = make([]string, 0, len(network.FConnections))
	for _, c := range network.FConnections {
		if conn.IsAmI(pCfg.GetAddress(), c) {
			continue
		}
		cfg.FConnections = append(cfg.FConnections, fmt.Sprintf("%s:%d", c.FHost, c.FPort))
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
	defaultNetwork := hiddenlake.GNetworks[hiddenlake.CDefaultNetwork]
	return &SConfig{
		FSettings: &SConfigSettings{
			FMessageSizeBytes: defaultNetwork.FMessageSizeBytes,
			FWorkSizeBits:     defaultNetwork.FWorkSizeBits,
			FMessagesCapacity: hlt_settings.CDefaultMessagesCapacity,
			FDatabaseEnabled:  hlt_settings.CDefaultDatabaseEnabled,
		},
		FLogging: []string{logger.CLogInfo, logger.CLogWarn, logger.CLogErro},
		FAddress: &SAddress{
			FTCP:   hlt_settings.CDefaultTCPAddress,
			FHTTP:  hlt_settings.CDefaultHTTPAddress,
			FPPROF: "",
		},
		FConnections: []string{
			hlt_settings.CDefaultConnectionAddress,
		},
		FConsumers: []string{},
	}
}
