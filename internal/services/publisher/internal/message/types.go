package message

import (
	"github.com/number571/hidden-lake/pkg/api/services/publisher/client/dto"
)

type IMessageContainer interface {
	GetFriend() string
	GetMessageUID() string
	GetMessageBody() dto.IMessageBody
}
