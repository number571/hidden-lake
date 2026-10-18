package dto

import (
	"encoding/json"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IMessage = &SMessage{}
)

type SMessage struct {
	FHead *SMessageHead `json:"head"`
	FBody *SMessageBody `json:"body"`
}

func NewMessage(pHead IMessageHead, pBody IMessageBody) IMessage {
	return &SMessage{
		FHead: pHead.(*SMessageHead),
		FBody: pBody.(*SMessageBody),
	}
}

func LoadMessage(pData interface{}) (IMessage, error) {
	var msgBytes []byte

	switch x := pData.(type) {
	case []byte:
		msgBytes = x
	case string:
		msgBytes = []byte(x)
	default:
		return nil, ErrUnknownType
	}

	msg := &SMessage{}
	if err := json.Unmarshal(msgBytes, msg); err != nil {
		return nil, err
	}

	if err := ValidateMessageHead(msg.GetHead()); err != nil {
		return nil, err
	}

	if err := ValidateMessageBody(msg.GetBody()); err != nil {
		return nil, err
	}

	return msg, nil
}

func (p *SMessage) GetHead() IMessageHead {
	return p.FHead
}

func (p *SMessage) GetBody() IMessageBody {
	return p.FBody
}

func (p *SMessage) ToBytes() []byte {
	return encoding.SerializeJSON(p)
}

func (p *SMessage) ToString() string {
	return string(p.ToBytes())
}
