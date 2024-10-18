package privkey

import (
	"os"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/utils"
)

func GetPrivKey(pKeyPath string) (asymmetric.IPrivKey, error) {
	if _, err := os.Stat(pKeyPath); os.IsNotExist(err) {
		privKey := asymmetric.NewPrivKey()
		if err := os.WriteFile(pKeyPath, []byte(privKey.ToString()), 0o600); err != nil {
			return nil, utils.MergeErrors(ErrWritePrivateKey, err)
		}
		return privKey, nil
	}
	privKeyStr, err := os.ReadFile(pKeyPath)
	if err != nil {
		return nil, utils.MergeErrors(ErrReadPrivateKey, err)
	}
	privKey := asymmetric.LoadPrivKey(string(privKeyStr))
	if privKey == nil {
		return nil, ErrInvalidPrivateKey
	}
	return privKey, nil
}
