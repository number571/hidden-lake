package database

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/types"
)

type IKVDatabase interface {
	types.ICloser

	Size(IRelation) uint64
	Push(IRelation, IMessage) error
	Load(IRelation, uint64, uint64) ([]IMessage, error)
}

type IRelation interface {
	IAm() asymmetric.IPubKeyChain
	Friend() asymmetric.IPubKeyChain
}

type IMessage interface {
	IsIncoming() bool
	GetTimestamp() string
	GetMessage() []byte
	ToBytes() []byte
}
