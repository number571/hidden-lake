package config

import (
	"os"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/utils"
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
	FNetworkKey       string `json:"network_key,omitempty" yaml:"network_key,omitempty"`
}

type SConfig struct {
	fFilepath string
	fMutex    sync.RWMutex
	fLogging  logger.ILogging
	fFriends  map[string]asymmetric.IPubKey

	FSettings    *SConfigSettings     `yaml:"settings"`
	FLogging     []string             `yaml:"logging,omitempty"`
	FAddress     *SAddress            `yaml:"address,omitempty"`
	FServices    map[string]*SService `yaml:"services,omitempty"`
	FConnections []string             `yaml:"connections,omitempty"`
	FFriends     map[string]string    `yaml:"friends,omitempty"`
}

type SService struct {
	FHost string `yaml:"host"`
}

type SAddress struct {
	FTCP   string `yaml:"tcp,omitempty"`
	FHTTP  string `yaml:"http,omitempty"`
	FPPROF string `yaml:"pprof,omitempty"`
}

func BuildConfig(pFilepath string, pCfg *SConfig) (IConfig, error) {
	if _, err := os.Stat(pFilepath); !os.IsNotExist(err) {
		return nil, utils.MergeErrors(ErrConfigAlreadyExist, err)
	}

	pCfg.fFilepath = pFilepath
	if err := pCfg.initConfig(); err != nil {
		return nil, utils.MergeErrors(ErrInitConfig, err)
	}

	if err := os.WriteFile(pFilepath, encoding.SerializeYAML(pCfg), 0o600); err != nil {
		return nil, utils.MergeErrors(ErrWriteConfig, err)
	}

	return pCfg, nil
}

func LoadConfig(pFilepath string) (IConfig, error) {
	if _, err := os.Stat(pFilepath); os.IsNotExist(err) {
		return nil, utils.MergeErrors(ErrConfigNotFound, err)
	}

	bytes, err := os.ReadFile(pFilepath)
	if err != nil {
		return nil, utils.MergeErrors(ErrReadConfig, err)
	}

	cfg := new(SConfig)
	if err := encoding.DeserializeYAML(bytes, cfg); err != nil {
		return nil, utils.MergeErrors(ErrDeserializeConfig, err)
	}

	cfg.fFilepath = pFilepath
	if err := cfg.initConfig(); err != nil {
		return nil, utils.MergeErrors(ErrInitConfig, err)
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

func (p *SConfigSettings) GetFetchTimeoutMS() uint64 {
	return p.FFetchTimeoutMS
}

func (p *SConfigSettings) GetQueuePeriodMS() uint64 {
	return p.FQueuePeriodMS
}

func (p *SConfigSettings) GetNetworkKey() string {
	return p.FNetworkKey
}

func (p *SConfig) GetSettings() IConfigSettings {
	return p.FSettings
}

func (p *SConfig) isValid() bool {
	for _, v := range p.FServices {
		if v.FHost == "" {
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

	for k, v := range p.FServices {
		if v == nil {
			p.FServices[k] = new(SService)
		}
	}

	if !p.isValid() {
		return ErrInvalidConfig
	}

	if err := p.loadPubKeys(); err != nil {
		return utils.MergeErrors(ErrLoadPublicKey, err)
	}

	if err := p.loadLogging(); err != nil {
		return utils.MergeErrors(ErrLoadLogging, err)
	}

	return nil
}

func (p *SConfig) loadLogging() error {
	result, err := logger.LoadLogging(p.FLogging)
	if err != nil {
		return utils.MergeErrors(ErrInvalidLogging, err)
	}
	p.fLogging = result
	return nil
}

func (p *SConfig) loadPubKeys() error {
	p.fFriends = make(map[string]asymmetric.IPubKey)

	mapping := make(map[string]struct{})
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

func (p *SConfig) GetConnections() []string {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	return p.FConnections
}

func (p *SConfig) GetService(name string) (IService, bool) {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	service, ok := p.FServices[name]
	return service, ok
}

func (p *SService) GetHost() string {
	return p.FHost
}

func (p *SAddress) GetTCP() string {
	return p.FTCP
}

func (p *SAddress) GetHTTP() string {
	return p.FHTTP
}

func (p *SAddress) GetPPROF() string {
	return p.FPPROF
}
