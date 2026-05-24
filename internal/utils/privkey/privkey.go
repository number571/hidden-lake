package privkey

import (
	"errors"
	"fmt"
	"os"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func GetPrivKey(pKeyPath string) (asymmetric.IPrivKey, error) {
	pubKeyPath := fmt.Sprintf("%s.pub", pKeyPath) //nolint:perfsprint
	if _, err := os.Stat(pKeyPath); os.IsNotExist(err) {
		privKey := asymmetric.NewPrivKey()
		if err := os.WriteFile(pKeyPath, []byte(privKey.ToString()), 0600); err != nil {
			return nil, errors.Join(ErrWritePrivateKey, err)
		}
		if err := os.WriteFile(pubKeyPath, []byte(privKey.GetPubKey().ToString()), 0600); err != nil {
			return nil, errors.Join(ErrWritePublicKey, err)
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
	if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) { //nolint:nestif
		if err := os.WriteFile(pubKeyPath, []byte(privKey.GetPubKey().ToString()), 0600); err != nil {
			return nil, errors.Join(ErrWritePublicKey, err)
		}
	} else {
		pubKeyStr, err := os.ReadFile(pubKeyPath) //nolint:gosec
		if err != nil {
			return nil, errors.Join(ErrReadPublicKey, err)
		}
		pubKey := asymmetric.LoadPubKey(string(pubKeyStr))
		if pubKey == nil {
			return nil, ErrInvalidPublicKey
		}
		if privKey.GetPubKey().GetHasher().ToString() != pubKey.GetHasher().ToString() {
			return nil, ErrNotLinkedPublicKeyToPrivate
		}
	}
	return privKey, nil
}
