package message

import (
	"github.com/number571/hidden-lake/pkg/api/services/notifier/client/dto"
)

type IMessageContainer interface {
	WithMessageUID(string) IMessageContainer

	GetFriend() string
	GetMessageUID() string
	GetMessageBody() dto.IMessageBody
}
