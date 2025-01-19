package config

import (
	"errors"
	"os"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

var (
	_ IConfigSettings = &SConfigSettings{}
	_ IConfig         = &SConfig{}
	_ IAddress        = &SAddress{}
)

type SConfigSettings struct {
	FMessageSizeBytes uint64 `json:"message_size_bytes" yaml:"message_size_bytes"`
	FFetchTimeoutMS   uint64 `json:"fetch_timeout_ms" yaml:"fetch_timeout_ms"`
	FQueuePeriodMS    uint64 `json:"queue_period_ms" yaml:"queue_period_ms"`
	FWorkSizeBits     uint64 `json:"work_size_bits,omitempty" yaml:"work_size_bits,omitempty"`
	FQBPConsumers     uint64 `json:"qbp_consumers,omitempty" yaml:"qbp_consumers,omitempty"`
	FPowParallel      uint64 `json:"pow_parallel,omitempty" yaml:"pow_parallel,omitempty"`
	FNetworkKey       string `json:"network_key,omitempty" yaml:"network_key,omitempty"`
}

type SConfig struct {
	fFilepath string
	fMutex    sync.RWMutex
	fLogging  logger.ILogging
	fFriends  map[string]asymmetric.IPubKey

	FSettings  *SConfigSettings  `yaml:"settings"`
	FLogging   []string          `yaml:"logging,omitempty"`
	FAddress   *SAddress         `yaml:"address,omitempty"`
	FServices  map[string]string `yaml:"services,omitempty"`
	FEndpoints []string          `yaml:"endpoints,omitempty"`
	FFriends   map[string]string `yaml:"friends,omitempty"`
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
		return nil, errors.Join(ErrConfigNotFound, err)
	}

	bytes, err := os.ReadFile(pFilepath)
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

func (p *SConfigSettings) GetMessageSizeBytes() uint64 {
	return p.FMessageSizeBytes
}

func (p *SConfigSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}

func (p *SConfigSettings) GetEncKeySizeBytes() uint64 {
	return asymmetric.CKEMCiphertextSize
}

func (p *SConfigSettings) GetFetchTimeout() time.Duration {
	return time.Duration(p.FFetchTimeoutMS) * time.Millisecond
}

func (p *SConfigSettings) GetQueuePeriod() time.Duration {
	return time.Duration(p.FQueuePeriodMS) * time.Millisecond
}

func (p *SConfigSettings) GetNetworkKey() string {
	return p.FNetworkKey
}

func (p *SConfigSettings) GetQBPConsumers() uint64 {
	return p.FQBPConsumers
}

func (p *SConfigSettings) GetPowParallel() uint64 {
	return p.FPowParallel
}

func (p *SConfig) GetSettings() IConfigSettings {
	return p.FSettings
}

func (p *SConfig) isValid() bool {
	for _, v := range p.FServices {
		if v == "" {
			return false
		}
	}
	return true &&
		p.FSettings.FMessageSizeBytes != 0 &&
		p.FSettings.FQueuePeriodMS != 0 &&
		p.FSettings.FFetchTimeoutMS != 0
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

	if err := p.loadPubKeys(); err != nil {
		return errors.Join(ErrLoadPublicKey, err)
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

func (p *SConfig) loadPubKeys() error {
	p.fFriends = make(map[string]asymmetric.IPubKey, len(p.FFriends))
	mapping := make(map[string]struct{}, len(p.FFriends))

	for name, val := range p.FFriends {
		if _, ok := mapping[val]; ok {
			return ErrDuplicatePublicKey
		}
		mapping[val] = struct{}{}

		pubKey := asymmetric.LoadPubKey(val)
		if pubKey == nil {
			return ErrInvalidPublicKey
		}
		p.fFriends[name] = pubKey
	}

	return nil
}

func (p *SConfig) GetFriends() map[string]asymmetric.IPubKey {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	result := make(map[string]asymmetric.IPubKey, len(p.FFriends))
	for k, v := range p.fFriends {
		result[k] = v
	}
	return result
}

func (p *SConfig) GetLogging() logger.ILogging {
	return p.fLogging
}

func (p *SConfig) GetAddress() IAddress {
	return p.FAddress
}

func (p *SConfig) GetEndpoints() []string {
	return p.FEndpoints
}

func (p *SConfig) GetService(name string) (string, bool) {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	service, ok := p.FServices[name]
	return service, ok
}

func (p *SAddress) GetExternal() string {
	return p.FExternal
}

func (p *SAddress) GetInternal() string {
	return p.FInternal
}
