package message

type sMessageContainer struct {
	fFriend  string
	fMessage IMessage
}

func newMessageContainer(pFriend string, pMessage IMessage) IMessageContainer {
	return &sMessageContainer{
		fFriend:  pFriend,
		fMessage: pMessage,
	}
}

func (p *sMessageContainer) GetFriend() string {
	return p.fFriend
}

func (p *sMessageContainer) GetMessage() IMessage {
	return p.fMessage
}
