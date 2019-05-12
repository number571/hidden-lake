package models

import (
	"crypto/rsa"
)

type Private struct {
	Key struct {
		P2P *rsa.PrivateKey
		F2F *rsa.PrivateKey
	}
	Data struct {
		P2P string
		F2F string
	}
}
