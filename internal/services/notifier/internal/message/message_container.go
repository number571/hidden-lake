package message

var (
	_ IMessageContainer = &sMessageContainer{}
)

type sMessageContainer struct {
	fFriend      string
	fMessageBody []byte
}

func NewMessageContainer(pFriend string, pBody []byte) IMessageContainer {
	return &sMessageContainer{
		fFriend:      pFriend,
		fMessageBody: pBody,
	}
}

func (p *sMessageContainer) GetFriend() string {
	return p.fFriend
}

func (p *sMessageContainer) GetMessageBody() []byte {
	return p.fMessageBody
}
