package message

import (
	"github.com/number571/hidden-lake/pkg/api/services/messenger/client/dto"
)

type IMessageBroker interface {
	Produce(string, dto.IMessage)
	Consume(string) <-chan IMessageContainer
}

type IMessageContainer interface {
	GetFriend() string
	GetMessage() dto.IMessage
}
