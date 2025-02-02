package config

import (
	"errors"
	"os"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
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

func (p *sEditor) UpdateFriends(pFriends map[string]asymmetric.IPubKey) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	filepath := p.fConfig.fFilepath
	icfg, err := LoadConfig(filepath)
	if err != nil {
		return errors.Join(ErrLoadConfig, err)
	}

	if hasDuplicatePubKeys(pFriends) {
		return ErrDuplicatePublicKey
	}

	cfg := icfg.(*SConfig)
	cfg.fFriends = pFriends
	cfg.FFriends = pubKeysToStrings(pFriends)
	if err := os.WriteFile(filepath, encoding.SerializeYAML(cfg), 0o600); err != nil {
		return errors.Join(ErrWriteConfig, err)
	}

	p.fConfig.fMutex.Lock()
	defer p.fConfig.fMutex.Unlock()

	p.fConfig.fFriends = cfg.fFriends
	p.fConfig.FFriends = cfg.FFriends
	return nil
}

func pubKeysToStrings(pPubKeys map[string]asymmetric.IPubKey) map[string]string {
	result := make(map[string]string, len(pPubKeys))
	for name, pubKey := range pPubKeys {
		result[name] = pubKey.ToString()
	}
	return result
}

func hasDuplicatePubKeys(pPubKeys map[string]asymmetric.IPubKey) bool {
	mapping := make(map[string]struct{}, len(pPubKeys))
	for _, pubKey := range pPubKeys {
		pubStr := pubKey.ToString()
		if _, ok := mapping[pubStr]; ok {
			return true
		}
		mapping[pubStr] = struct{}{}
	}
	return false
}

// func deleteDuplicateStrings(pStrs []string) []string {
// 	result := make([]string, 0, len(pStrs))
// 	mapping := make(map[string]struct{}, len(pStrs))
// 	for _, s := range pStrs {
// 		if _, ok := mapping[s]; ok {
// 			continue
// 		}
// 		mapping[s] = struct{}{}
// 		result = append(result, s)
// 	}
// 	return result
// }
