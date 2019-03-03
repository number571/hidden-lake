package crypto

import (
	"crypto/rsa"
	"crypto/rand"
	"crypto/sha256"
)

func EncryptRSA(pub *rsa.PublicKey, data []byte) ([]byte, error) {
    return rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, data, nil)
}
