package dto

import (
	"encoding/json"
	"time"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IMessage = &SMessage{}
)

type SMessage struct {
	FIncoming  bool   `json:"incoming"`
	FMessage   string `json:"message"`
	FTimestamp uint64 `json:"timestamp"`
}

func NewMessage(pIsIncoming bool, pMessage string, pTime time.Time) IMessage {
	return &SMessage{
		FIncoming:  pIsIncoming,
		FMessage:   pMessage,
		FTimestamp: uint64(pTime.Unix()), //nolint:gosec
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
	return msg, nil
}

func (p *SMessage) IsIncoming() bool {
	return p.FIncoming
}

func (p *SMessage) GetMessage() string {
	return p.FMessage
}

func (p *SMessage) GetTimestamp() string {
	return time.Unix(int64(p.FTimestamp), 0).Format(time.DateTime) //nolint:gosec
}

func (p *SMessage) ToBytes() []byte {
	return encoding.SerializeJSON(p)
}

func (p *SMessage) ToString() string {
	return string(p.ToBytes())
}
