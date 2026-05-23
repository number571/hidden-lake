package config

import (
	"errors"
	"os"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/encoding"
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

func (p *sEditor) UpdateFriends(pFriends map[string]layer2.IParticipantKey) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	filepath := p.fConfig.fFilepath
	icfg, err := LoadConfig(filepath)
	if err != nil {
		return errors.Join(ErrLoadConfig, err)
	}

	if hasDuplicateParticipantKeys(pFriends) {
		return ErrDuplicateParticipantKey
	}

	cfg := icfg.(*SConfig)
	cfg.fFriends = pFriends
	cfg.FFriends = participantKeysToStrings(pFriends)
	if err := os.WriteFile(filepath, encoding.SerializeYAML(cfg), 0600); err != nil {
		return errors.Join(ErrWriteConfig, err)
	}

	p.fConfig.fMutex.Lock()
	defer p.fConfig.fMutex.Unlock()

	p.fConfig.fFriends = cfg.fFriends
	p.fConfig.FFriends = cfg.FFriends
	return nil
}

func participantKeysToStrings(pFriends map[string]layer2.IParticipantKey) map[string]string {
	result := make(map[string]string, len(pFriends))
	for name, pubKey := range pFriends {
		result[name] = pubKey.ToString()
	}
	return result
}

func hasDuplicateParticipantKeys(pFriends map[string]layer2.IParticipantKey) bool {
	mapping := make(map[string]struct{}, len(pFriends))
	for _, pubKey := range pFriends {
		pubStr := pubKey.ToString()
		if _, ok := mapping[pubStr]; ok {
			return true
		}
		mapping[pubStr] = struct{}{}
	}
	return false
}
