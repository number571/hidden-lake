package config

import (
	"errors"
	"os"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

var (
	_ IConfig  = &SConfig{}
	_ IAddress = &SAddress{}
)

type SConfigSettings struct {
	FMessageSizeBytes uint64 `json:"message_size_bytes" yaml:"message_size_bytes"`
	FWorkSizeBits     uint64 `json:"work_size_bits,omitempty" yaml:"work_size_bits,omitempty"`
	FNetworkKey       string `json:"network_key,omitempty" yaml:"network_key,omitempty"`
}

type SConfig struct {
	fFilepath string
	fMutex    sync.RWMutex
	fFriends  map[string]asymmetric.IPubKey
	fLogging  logger.ILogging

	FSettings *SConfigSettings  `yaml:"settings"`
	FLogging  []string          `yaml:"logging,omitempty"`
	FAddress  *SAddress         `yaml:"address"`
	FFriends  map[string]string `yaml:"friends,omitempty"`
}

type SAddress struct {
	FHTTP  string `yaml:"http"`
	FPPROF string `yaml:"pprof,omitempty"`
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

func (p *SConfig) isValid() bool {
	return true &&
		p.FSettings.FMessageSizeBytes != 0 &&
		p.FAddress.FHTTP != ""
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

func (p *SConfig) GetSettings() IConfigSettings {
	return p.FSettings
}

func (p *SConfig) GetAddress() IAddress {
	return p.FAddress
}

func (p *SAddress) GetHTTP() string {
	return p.FHTTP
}

func (p *SAddress) GetPPROF() string {
	return p.FPPROF
}

func (p *SConfig) GetFriends() map[string]asymmetric.IPubKey {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	result := make(map[string]asymmetric.IPubKey, len(p.fFriends))
	for k, v := range p.fFriends {
		result[k] = v
	}
	return result
}

func (p *SConfigSettings) GetNetworkKey() string {
	return p.FNetworkKey
}

func (p *SConfigSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}

func (p *SConfigSettings) GetMessageSizeBytes() uint64 {
	return p.FMessageSizeBytes
}

func (p *SConfigSettings) GetEncKeySizeBytes() uint64 {
	return asymmetric.CKEMCiphertextSize
}

func (p *SConfig) GetLogging() logger.ILogging {
	return p.fLogging
}
