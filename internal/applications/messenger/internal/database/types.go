package database

import (
	"io"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

type IKVDatabase interface {
	io.Closer

	Size(IRelation) uint64
	Push(IRelation, IMessage) error
	Load(IRelation, uint64, uint64) ([]IMessage, error)
}

type IRelation interface {
	IAm() asymmetric.IPubKey
	Friend() asymmetric.IPubKey
}

type IMessage interface {
	IsIncoming() bool
	GetTimestamp() string
	GetMessage() []byte
	ToBytes() []byte
}
