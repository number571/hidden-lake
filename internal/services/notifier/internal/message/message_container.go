package message

import "github.com/number571/hidden-lake/pkg/api/services/notifier/client/dto"

var (
	_ IMessageContainer = &sMessageContainer{}
)

type sMessageContainer struct {
	fFriend      string
	fMessageUID  string
	fMessageBody dto.IMessageBody
}

func NewMessageContainer(pFriend string, pMessageBody dto.IMessageBody) IMessageContainer {
	return &sMessageContainer{
		fFriend:      pFriend,
		fMessageBody: pMessageBody,
	}
}

func (p *sMessageContainer) WithMessageUID(pUID string) IMessageContainer {
	p.fMessageUID = pUID
	return p
}

func (p *sMessageContainer) GetFriend() string {
	return p.fFriend
}

func (p *sMessageContainer) GetMessageUID() string {
	return p.fMessageUID
}

func (p *sMessageContainer) GetMessageBody() dto.IMessageBody {
	return p.fMessageBody
}
