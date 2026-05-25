package keys

import (
	"errors"
	"os"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func GetPubKey(pPrivKey asymmetric.IPrivKey, pPubKeyPath string) (asymmetric.IPubKey, error) {
	if _, err := os.Stat(pPubKeyPath); os.IsNotExist(err) { //nolint:nestif
		pubKey := pPrivKey.GetPubKey()
		if err := os.WriteFile(pPubKeyPath, []byte(pubKey.ToString()), 0600); err != nil {
			return nil, errors.Join(ErrWritePublicKey, err)
		}
		return pubKey, nil
	}
	pubKeyStr, err := os.ReadFile(pPubKeyPath) //nolint:gosec
	if err != nil {
		return nil, errors.Join(ErrReadPublicKey, err)
	}
	pubKey := asymmetric.LoadPubKey(string(pubKeyStr))
	if pubKey == nil {
		return nil, ErrInvalidPublicKey
	}
	if pPrivKey.GetPubKey().GetHasher().ToString() != pubKey.GetHasher().ToString() {
		return nil, ErrNotLinkedPublicKeyToPrivate
	}
	return pubKey, nil
}
