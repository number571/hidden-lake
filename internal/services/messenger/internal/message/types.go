package message

import (
	"github.com/number571/hidden-lake/pkg/api/services/messenger/client/dto"
)

type IMessageContainer interface {
	GetFriend() string
	GetMessage() dto.IMessage
}
