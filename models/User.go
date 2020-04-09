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
	ChatMap  ChatMap
}

type ChatMap struct {
	Owner  map[string]bool // list of users in own chat
	Member map[string]bool // list of global chats by founders hashnames
}

type Auth struct {
	Hashpasw string // hash(hash(password+salt))
	Pasw     []byte // hash(password+salt)
	Salt     string // base64(random_bytes(16))
}

type Keys struct {
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

type Session struct {
	Option ChatOption
	Socket *websocket.Conn
	Time   string
}

type ChatOption uint8

const (
	PRIVATE_OPTION ChatOption = 1
	GROUP_OPTION   ChatOption = 2
)
