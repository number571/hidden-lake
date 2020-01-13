package models

import (
	"crypto/rsa"
)

type Client struct {
	Hashname string
	Address  string
	Public   *rsa.PublicKey
}
