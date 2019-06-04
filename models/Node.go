package models

import (
	"net"
	"crypto/rsa"
)

type SessionKey struct {
	P2P map[string][]byte
	F2F map[string][]byte
}

type Address struct {
	P2P map[string]string
	F2F map[string]string
	C_S map[string]net.Conn
}

type ConnServer struct {
    Addr net.Conn
    Hash string
}

type Node struct {
	Address Address
	SessionKey SessionKey
	ConnServer ConnServer
	ConnectionMode map[string]ModeConn
	PublicKey map[string]*rsa.PublicKey
}
