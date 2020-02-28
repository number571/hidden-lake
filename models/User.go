package models

import (
	"crypto/rsa"
	"golang.org/x/net/websocket"
)

type User struct {
	Username string 
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
	Hashpasw string // hash(hash(password+salt))
	Pasw     []byte // hash(password+salt)
	Salt     string // base64(random_bytes(8))
}

type Keys struct {
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

type Session struct {
	Socket *websocket.Conn
	Time   string
}
