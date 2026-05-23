package database

import (
	"io"

	message "github.com/number571/hidden-lake/pkg/api/services/messenger/client/dto"
)

type IKVDatabase interface {
	io.Closer

	Size(string) uint64
	Push(string, message.IMessage) error
	Load(string, uint64, uint64) ([]message.IMessage, error)
}
