//go:build !symmetric

package app

import (
	"path/filepath"

	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2/hybrid"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/privkey"
)

func getScheme(inputPath string, msgSizeBytes uint64) (layer2.IScheme, error) {
	keyPath := filepath.Join(inputPath, hlk_settings.CPathKey)
	privKey, err := privkey.GetPrivKey(keyPath)
	if err != nil {
		return nil, err
	}
	return hybrid.NewScheme(privKey, msgSizeBytes)
}
