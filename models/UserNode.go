package models

import (
    "crypto/rsa"
)

type Hash struct {
	P2P string
	F2F string
}

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

type Transportation struct {
    IPv4 string
    Port string
}

type UserNode struct {
    Mode ModeNet
    Hash Hash
    Public Public
    Private Private
    Connection
    Authorization
    Transportation
}
