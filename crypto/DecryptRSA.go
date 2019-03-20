package crypto

import (
	"crypto/rsa"
	"crypto/rand"
	"crypto/sha256"
)

func DecryptRSA(data []byte, priv *rsa.PrivateKey) ([]byte, error) {
    return rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, data, nil)
}
