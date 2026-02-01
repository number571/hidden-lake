package config

import (
	"errors"
	"os"
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/hidden-lake/internal/utils/slices"
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

func (p *sEditor) UpdateConnections(pConns []string) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	filepath := p.fConfig.fFilepath
	icfg, err := LoadConfig(filepath)
	if err != nil {
		return errors.Join(ErrLoadConfig, err)
	}

	cfg := icfg.(*SConfig)
	cfg.FConnections = slices.DeleteDuplicates(pConns)
	if err := os.WriteFile(filepath, encoding.SerializeYAML(cfg), 0600); err != nil {
		return errors.Join(ErrWriteConfig, err)
	}

	p.fConfig.fMutex.Lock()
	defer p.fConfig.fMutex.Unlock()

	p.fConfig.FConnections = cfg.FConnections
	return nil
}
