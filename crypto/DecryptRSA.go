package crypto

import (
	"crypto/rsa"
	"crypto/rand"
	"crypto/sha256"
)

func DecryptRSA(priv *rsa.PrivateKey, data []byte) ([]byte, error) {
    return rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, data, nil)
}
