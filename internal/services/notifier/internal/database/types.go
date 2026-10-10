package database

import (
	"io"

	message "github.com/number571/hidden-lake/pkg/api/services/notifier/client/dto"
)

type IKVDatabase interface {
	io.Closer

	Size(string) uint64
	Push(string, message.IMessage) (uint64, error)

	Recv(string, string) error

	LoadByUID(string, string) (message.IMessage, error)
	LoadByEnum(string, uint64) (message.IMessage, error)
}
