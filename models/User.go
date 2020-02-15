package models

import (
	"crypto/rsa"
	"golang.org/x/net/websocket"
)

type User struct {
	Hashname string // hash(pubkey)
	UsedF2F  bool
	Auth     Auth
	Keys     Keys
	Temp     Temp
	Session  Session
}

type Temp struct {
	FileList []File
	ConnList []Connect
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
