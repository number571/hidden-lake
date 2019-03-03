package encoding

import (
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "../utils"
)

func DecodePrivate(data string) *rsa.PrivateKey {
    block, _ := pem.Decode([]byte(data))

    priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    utils.CheckError(err)

    return priv
}
