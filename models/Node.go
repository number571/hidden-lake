package models

import (
	"crypto/rsa"
)

type SessionKey struct {
	P2P map[string][]byte
	F2F map[string][]byte
}

type Address struct {
	P2P map[string]string
	F2F map[string]string
}

type Node struct {
	Login map[string]string
	Connection map[string]int8
	PublicKey map[string]*rsa.PublicKey
	SessionKey SessionKey
	Address Address
}
