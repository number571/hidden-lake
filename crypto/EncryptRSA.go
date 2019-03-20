package crypto

import (
	"crypto/rsa"
	"crypto/rand"
	"crypto/sha256"
)

func EncryptRSA(data []byte, pub *rsa.PublicKey) ([]byte, error) {
    return rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, data, nil)
}
