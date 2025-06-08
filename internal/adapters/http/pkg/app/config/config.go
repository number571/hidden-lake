package config

import (
	"errors"
	"os"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/encoding"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

var (
	_ IConfig  = &SConfig{}
	_ IAddress = &SAddress{}
)

type SConfigSettings struct {
	FMessageSizeBytes uint64 `json:"message_size_bytes,omitempty" yaml:"message_size_bytes,omitempty"`
	FWorkSizeBits     uint64 `json:"work_size_bits,omitempty" yaml:"work_size_bits,omitempty"`
	FNetworkKey       string `json:"network_key,omitempty" yaml:"network_key,omitempty"`
	FDatabaseEnabled  bool   `json:"database_enabled,omitempty" yaml:"database_enabled,omitempty"`
	FSendTimeoutMS    uint64 `json:"send_timeout_ms,omitempty" yaml:"send_timeout_ms,omitempty"`
	FRecvTimeoutMS    uint64 `json:"recv_timeout_ms,omitempty" yaml:"recv_timeout_ms,omitempty"`
}

type SConfig struct {
	fFilepath string
	fMutex    sync.RWMutex
	fLogging  logger.ILogging

	FSettings    *SConfigSettings `yaml:"settings,omitempty"`
	FLogging     []string         `yaml:"logging,omitempty"`
	FAddress     *SAddress        `yaml:"address,omitempty"`
	FEndpoints   []string         `yaml:"endpoints,omitempty"`
	FConnections []string         `yaml:"connections,omitempty"`
}

type SAddress struct {
	FExternal string `yaml:"external,omitempty"`
	FInternal string `yaml:"internal,omitempty"`
}

func BuildConfig(pFilepath string, pCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(pFilepath); !os.IsNotExist(err) {
		return nil, errors.Join(ErrConfigAlreadyExist, err)
	}

	pCfg.fFilepath = pFilepath
	if err := pCfg.initConfig(); err != nil {
		return nil, errors.Join(ErrInitConfig, err)
	}

	if err := os.WriteFile(pFilepath, encoding.SerializeYAML(pCfg), 0o600); err != nil {
		return nil, errors.Join(ErrWriteConfig, err)
	}

	return pCfg, nil
}

func LoadConfig(pFilepath string) (IConfig, error) {
	if _, err := os.Stat(pFilepath); os.IsNotExist(err) {
		return nil, errors.Join(ErrConfigNotExist, err)
	}

	bytes, err := os.ReadFile(pFilepath) //nolint:gosec
	if err != nil {
		return nil, errors.Join(ErrReadConfig, err)
	}

	cfg := new(SConfig)
	if err := encoding.DeserializeYAML(bytes, cfg); err != nil {
		return nil, errors.Join(ErrDeserializeConfig, err)
	}

	cfg.fFilepath = pFilepath
	if err := cfg.initConfig(); err != nil {
		return nil, errors.Join(ErrInitConfig, err)
	}

	return cfg, nil
}

func (p *SConfig) isValid() bool {
	return true
}

func (p *SConfig) initConfig() error {
	if p.FSettings == nil {
		p.FSettings = new(SConfigSettings)
	}

	if p.FAddress == nil {
		p.FAddress = new(SAddress)
	}

	if !p.isValid() {
		return ErrInvalidConfig
	}

	if err := p.loadLogging(); err != nil {
		return errors.Join(ErrLoadLogging, err)
	}

	return nil
}

func (p *SConfig) loadLogging() error {
	result, err := logger.LoadLogging(p.FLogging)
	if err != nil {
		return errors.Join(ErrInvalidLogging, err)
	}
	p.fLogging = result
	return nil
}

func (p *SConfig) GetSettings() IConfigSettings {
	return p.FSettings
}

func (p *SConfig) GetLogging() logger.ILogging {
	return p.fLogging
}

func (p *SConfig) GetEndpoints() []string {
	return p.FEndpoints
}

func (p *SConfig) GetConnections() []string {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	return p.FConnections
}

func (p *SConfigSettings) GetMessageSizeBytes() uint64 {
	return p.FMessageSizeBytes
}

func (p *SConfigSettings) GetNetworkKey() string {
	return p.FNetworkKey
}

func (p *SConfigSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}

func (p *SConfigSettings) GetSendTimeout() time.Duration {
	return time.Duration(p.FSendTimeoutMS) * time.Millisecond // nolint: gosec
}

func (p *SConfigSettings) GetRecvTimeout() time.Duration {
	return time.Duration(p.FRecvTimeoutMS) * time.Millisecond // nolint: gosec
}

func (p *SConfigSettings) GetDatabaseEnabled() bool {
	return p.FDatabaseEnabled
}

func (p *SConfig) GetAddress() IAddress {
	return p.FAddress
}

func (p *SAddress) GetExternal() string {
	return p.FExternal
}

func (p *SAddress) GetInternal() string {
	return p.FInternal
}
