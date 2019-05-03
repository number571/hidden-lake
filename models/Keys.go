package models

import (
	"crypto/rsa"
)

type KeysF2F struct {
	NodeSessionKeyF2F map[string][]byte
}

type KeysP2P struct {
	NodePublicKey map[string]*rsa.PublicKey
	NodeSessionKey map[string][]byte
	KeysF2F
}

type Keys struct {
	PrivateKey *rsa.PrivateKey
	PublicKey *rsa.PublicKey
	PrivateData string
	PublicData string
}
