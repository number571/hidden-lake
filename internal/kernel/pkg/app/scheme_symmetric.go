//go:build symmetric

package app

import (
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2/symmetric"
)

func getScheme(_ string, msgSizeBytes uint64) (layer2.IScheme, error) {
	return symmetric.NewScheme(msgSizeBytes), nil
}
