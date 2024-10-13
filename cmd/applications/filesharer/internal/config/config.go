package config

import (
	"os"
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/utils"
	"github.com/number571/hidden-lake/internal/language"
	logger "github.com/number571/hidden-lake/internal/logger/std"
)

var (
	_ IConfigSettings = &SConfigSettings{}
	_ IConfig         = &SConfig{}
	_ IAddress        = &SAddress{}
)

type SConfigSettings struct {
	fMutex    sync.RWMutex
	fLanguage language.ILanguage

	FPageOffset uint64 `json:"page_offset" yaml:"page_offset"`
	FRetryNum   uint64 `json:"retry_num,omitempty" yaml:"retry_num,omitempty"`
	FLanguage   string `json:"language,omitempty" yaml:"language,omitempty"`
}

type SConfig struct {
	fFilepath string
	fLogging  logger.ILogging

	FSettings   *SConfigSettings `yaml:"settings"`
	FLogging    []string         `yaml:"logging,omitempty"`
	FAddress    *SAddress        `yaml:"address"`
	FConnection string           `yaml:"connection"`
}

type SAddress struct {
	FInterface string `yaml:"interface"`
	FIncoming  string `yaml:"incoming"`
	FPPROF     string `yaml:"pprof,omitempty"`
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
		return nil, utils.MergeErrors(ErrConfigNotExist, err)
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

func (p *SConfig) GetSettings() IConfigSettings {
	return p.FSettings
}

func (p *SConfigSettings) GetRetryNum() uint64 {
	return p.FRetryNum
}

func (p *SConfigSettings) GetPageOffset() uint64 {
	return p.FPageOffset
}

func (p *SConfigSettings) GetLanguage() language.ILanguage {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	return p.fLanguage
}

func (p *SConfigSettings) loadLanguage() error {
	res, err := language.ToILanguage(p.FLanguage)
	if err != nil {
		return utils.MergeErrors(ErrToLanguage, err)
	}
	p.fLanguage = res
	return nil
}

func (p *SConfig) isValid() bool {
	return true &&
		p.FConnection != "" &&
		p.FAddress.FInterface != "" &&
		p.FAddress.FIncoming != "" &&
		p.FSettings.FPageOffset != 0
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
		return utils.MergeErrors(ErrLoadLogging, err)
	}

	if err := p.FSettings.loadLanguage(); err != nil {
		return utils.MergeErrors(ErrLoadLanguage, err)
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

func (p *SConfig) GetAddress() IAddress {
	return p.FAddress
}

func (p *SConfig) GetConnection() string {
	return p.FConnection
}

func (p *SAddress) GetInterface() string {
	return p.FInterface
}

func (p *SAddress) GetIncoming() string {
	return p.FIncoming
}

func (p *SAddress) GetPPROF() string {
	return p.FPPROF
}

func (p *SConfig) GetLogging() logger.ILogging {
	return p.fLogging
}