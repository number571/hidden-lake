package crypto

import (
	"crypto/rsa"
	"crypto/rand"
	"../utils"
)

func GenerateKeys(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	utils.CheckError(err)
	return priv, &priv.PublicKey
}
