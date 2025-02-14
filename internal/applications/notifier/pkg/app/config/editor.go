package config

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/hidden-lake/internal/utils/language"
)

var (
	_ IEditor = &sEditor{}
)

type sEditor struct {
	fMutex  sync.Mutex
	fConfig *SConfig
}

func newEditor(pCfg IConfig) IEditor {
	if pCfg == nil {
		panic("cfg = nil")
	}
	v, ok := pCfg.(*SConfig)
	if !ok {
		panic("cfg is invalid")
	}
	return &sEditor{
		fConfig: v,
	}
}

func (p *sEditor) UpdateLanguage(pLang language.ILanguage) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	filepath := p.fConfig.fFilepath
	icfg, err := LoadConfig(filepath)
	if err != nil {
		return fmt.Errorf("load config (update language): %w", err)
	}

	cfg := icfg.(*SConfig)
	cfg.FSettings.FLanguage = language.FromILanguage(pLang)
	if err := os.WriteFile(filepath, encoding.SerializeYAML(cfg), 0o600); err != nil {
		return fmt.Errorf("write config (update language): %w", err)
	}

	p.fConfig.FSettings.fMutex.Lock()
	defer p.fConfig.FSettings.fMutex.Unlock()

	p.fConfig.FSettings.fLanguage = pLang
	return nil
}

func (p *sEditor) UpdateChannels(pChannels []string) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	filepath := p.fConfig.fFilepath
	icfg, err := LoadConfig(filepath)
	if err != nil {
		return errors.Join(ErrLoadConfig, err)
	}

	cfg := icfg.(*SConfig)
	cfg.FChannels = deleteDuplicateStrings(pChannels)
	if err := os.WriteFile(filepath, encoding.SerializeYAML(cfg), 0o600); err != nil {
		return errors.Join(ErrWriteConfig, err)
	}

	p.fConfig.fMutex.Lock()
	defer p.fConfig.fMutex.Unlock()

	p.fConfig.FChannels = cfg.FChannels
	return nil
}

func deleteDuplicateStrings(pStrs []string) []string {
	result := make([]string, 0, len(pStrs))
	mapping := make(map[string]struct{}, len(pStrs))
	for _, s := range pStrs {
		if _, ok := mapping[s]; ok {
			continue
		}
		mapping[s] = struct{}{}
		result = append(result, s)
	}
	return result
}
