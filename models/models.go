package models

import (
	"crypto/rsa"
)

type From struct {
	Address string
	Name string
}

type Head struct {
	Header string
	Mode string
}

type Transportation struct {
	Name string
	IPv4 string
	Port string
	Info string
}

type Keys struct {
	PrivateKey *rsa.PrivateKey
	PublicKey *rsa.PublicKey
	PrivateData string
	PublicData string
	NodePublicKey map[string]*rsa.PublicKey
	NodeSessionKey map[string][]byte
	NodeConnection map[string]int8
}

type Connection struct {
	TempConnect string
	TempProfile []string
	TempArchive []string
	Connections []string
}

type ChatMessages struct {
	GlobalMessages []string
	LocalMessages map[string][]string
}
