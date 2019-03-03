package encoding

import (
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
)

func EncodePublic(pub *rsa.PublicKey) []byte {
    return pem.EncodeToMemory(
        &pem.Block{
            Type: "RSA PUBLIC KEY",
            Bytes: x509.MarshalPKCS1PublicKey(pub),
        },
    )
}
