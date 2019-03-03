package encoding

import (
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "../utils"
)

func DecodePublic(data string) *rsa.PublicKey {
    block, _ := pem.Decode([]byte(data))

    pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
    utils.CheckError(err)

    return pub
}
