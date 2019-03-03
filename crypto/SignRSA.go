package crypto

import (
	"crypto"
	"crypto/rsa"
	"crypto/rand"
)

func SignRSA(priv *rsa.PrivateKey, data []byte) ([]byte, error) {
    return rsa.SignPSS(rand.Reader, priv, crypto.SHA256, HashSum(data), nil)
}
