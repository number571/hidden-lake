package message

import "github.com/number571/hidden-lake/pkg/api/services/messenger/client/dto"

type sMessageContainer struct {
	fFriend  string
	fMessage dto.IMessage
}

func newMessageContainer(pFriend string, pMessage dto.IMessage) IMessageContainer {
	return &sMessageContainer{
		fFriend:  pFriend,
		fMessage: pMessage,
	}
}

func (p *sMessageContainer) GetFriend() string {
	return p.fFriend
}

func (p *sMessageContainer) GetMessage() dto.IMessage {
	return p.fMessage
}
