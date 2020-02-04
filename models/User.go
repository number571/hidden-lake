package models

import (
	"crypto/rsa"
	"golang.org/x/net/websocket"
)

type User struct {
	Hashname string // hash(pubkey)
	Auth     Auth
	Keys     Keys
	FileList []File
	Session  Session
}

type Auth struct {
	Hashpasw string // hash(hash(username+password))
	Pasw     []byte // hash(username+password)
}

type Keys struct {
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

type Session struct {
	Socket *websocket.Conn
	Time   string
}
