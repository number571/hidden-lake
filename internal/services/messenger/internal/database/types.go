package database

import (
	"io"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/client/message"
)

type IKVDatabase interface {
	io.Closer

	Size(IRelation) uint64
	Push(IRelation, message.IMessage) error
	Load(IRelation, uint64, uint64) ([]message.IMessage, error)
}

type IRelation interface {
	IAm() asymmetric.IPubKey
	Friend() asymmetric.IPubKey
}
