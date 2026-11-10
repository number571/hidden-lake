package message

import (
	message_dto "github.com/number571/hidden-lake/pkg/api/services/publisher/client/dto"
)

var (
	_ IMessageContainer = &sMessageContainer{}
)

type sMessageContainer struct {
	fFriend string
	fNID    uint64
	fMsg    message_dto.IMessage
}

func NewMessageContainer(pFriend string, pNID uint64, pMsg message_dto.IMessage) IMessageContainer {
	return &sMessageContainer{
		fFriend: pFriend,
		fNID:    pNID,
		fMsg:    pMsg,
	}
}

func (p *sMessageContainer) GetFriend() string {
	return p.fFriend
}

func (p *sMessageContainer) GetMsgNID() uint64 {
	return p.fNID
}

func (p *sMessageContainer) GetMessage() message_dto.IMessage {
	return p.fMsg
}
