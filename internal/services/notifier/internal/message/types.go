package message

import (
	"github.com/number571/hidden-lake/pkg/api/services/notifier/client/dto"
)

type IMessageContainer interface {
	GetFriend() string
	GetMessage() dto.IMessage
}
