package encoding

import (
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
)

func EncodePrivate(priv *rsa.PrivateKey) []byte {
    return pem.EncodeToMemory(
        &pem.Block{
            Type: "RSA PRIVATE KEY",
            Bytes: x509.MarshalPKCS1PrivateKey(priv),
        },
    )
}
