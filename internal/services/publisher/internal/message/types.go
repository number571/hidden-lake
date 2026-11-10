package message

import message_dto "github.com/number571/hidden-lake/pkg/api/services/publisher/client/dto"

type IMessageContainer interface {
	GetFriend() string
	GetMsgNID() uint64
	GetMessage() message_dto.IMessage
}
