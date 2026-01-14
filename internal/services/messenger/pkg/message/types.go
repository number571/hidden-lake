package message

import "github.com/number571/go-peer/pkg/types"

type IMessageBroker interface {
	Produce(string, IMessage)
	Consume(string) chan IMessageContainer
}

type IMessageContainer interface {
	GetFriend() string
	GetMessage() IMessage
}

type IMessage interface {
	types.IConverter

	IsIncoming() bool
	GetTimestamp() string
	GetMessage() string
}
