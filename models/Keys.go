package models

import (
	"crypto/rsa"
)

type Keys struct {
	PrivateKey *rsa.PrivateKey
	PublicKey *rsa.PublicKey
	PrivateData string
	PublicData string
	NodePublicKey map[string]*rsa.PublicKey
	NodeSessionKey map[string][]byte
	NodeConnection map[string]int8
}
