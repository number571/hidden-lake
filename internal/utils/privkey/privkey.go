package privkey

import (
	"errors"
	"os"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func GetPrivKey(pKeyPath string) (asymmetric.IPrivKey, error) {
	if _, err := os.Stat(pKeyPath); os.IsNotExist(err) {
		privKey := asymmetric.NewPrivKey()
		if err := os.WriteFile(pKeyPath, []byte(privKey.ToString()), 0600); err != nil {
			return nil, errors.Join(ErrWritePrivateKey, err)
		}
		return privKey, nil
	}
	privKeyStr, err := os.ReadFile(pKeyPath) //nolint:gosec
	if err != nil {
		return nil, errors.Join(ErrReadPrivateKey, err)
	}
	privKey := asymmetric.LoadPrivKey(string(privKeyStr))
	if privKey == nil {
		return nil, ErrInvalidPrivateKey
	}
	return privKey, nil
}
