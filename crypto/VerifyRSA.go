package crypto

import (
	"crypto"
	"crypto/rsa"
)

func VerifyRSA(pub *rsa.PublicKey, data, sign []byte) error {
    return rsa.VerifyPSS(pub, crypto.SHA256, HashSum(data), sign, nil)
}
