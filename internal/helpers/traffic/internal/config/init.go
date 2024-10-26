package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/utils"
	hiddenlake "github.com/number571/hidden-lake"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"

	hlt_settings "github.com/number571/hidden-lake/internal/helpers/traffic/pkg/settings"
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

	cfg.FSettings.FMessageSizeBytes = network.FMessageSizeBytes
	cfg.FSettings.FWorkSizeBits = network.FWorkSizeBits
	cfg.FSettings.FNetworkKey = pUseNetwork

	cfg.FConnections = make([]string, 0, len(network.FConnections))
	for _, c := range network.FConnections {
		if isAmI(pCfg, c) {
			continue
		}
		cfg.FConnections = append(cfg.FConnections, fmt.Sprintf("%s:%d", c.FHost, c.FPort))
	}

	if err := os.WriteFile(cfg.fFilepath, encoding.SerializeYAML(cfg), 0o600); err != nil {
		return nil, utils.MergeErrors(ErrRebuildConfig, ErrWriteConfig, err)
	}

	rCfg, err := LoadConfig(cfg.fFilepath)
	if err != nil {
		return nil, utils.MergeErrors(ErrRebuildConfig, ErrLoadConfig, err)
	}

	return rCfg, nil
}

func isAmI(pCfg IConfig, conn hiddenlake.SConnection) bool {
	splited := strings.Split(pCfg.GetAddress().GetTCP(), ":")
	if len(splited) < 2 {
		return false
	}
	tcpPort, _ := strconv.Atoi(splited[1])
	if conn.FHost == "localhost" || conn.FHost == "127.0.0.1" {
		if conn.FPort == uint16(tcpPort) {
			return true
		}
	}
	return false
}

func initConfig() *SConfig {
	defaultNetwork, ok := hiddenlake.GNetworks[hiddenlake.CDefaultNetwork]
	if !ok {
		panic("get default network")
	}
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
