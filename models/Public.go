package models

import (
	"crypto/rsa"
)

type Public struct {
	Key struct {
		P2P *rsa.PublicKey
		F2F *rsa.PublicKey
	}
	Data struct {
		P2P string
		F2F string
	}
}
