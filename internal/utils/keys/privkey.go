package keys

import (
	"errors"
	"os"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func GetPrivKey(pPrivKeyPath string) (asymmetric.IPrivKey, error) {
	if _, err := os.Stat(pPrivKeyPath); os.IsNotExist(err) {
		privKey := asymmetric.NewPrivKey()
		if err := os.WriteFile(pPrivKeyPath, []byte(privKey.ToString()), 0600); err != nil {
			return nil, errors.Join(ErrWritePrivateKey, err)
		}
		return privKey, nil
	}
	privKeyStr, err := os.ReadFile(pPrivKeyPath) //nolint:gosec
	if err != nil {
		return nil, errors.Join(ErrReadPrivateKey, err)
	}
	privKey := asymmetric.LoadPrivKey(string(privKeyStr))
	if privKey == nil {
		return nil, ErrInvalidPrivateKey
	}
	return privKey, nil
}
