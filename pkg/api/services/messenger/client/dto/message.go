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

func NewMessage(pIsIncoming bool, pMessage string) IMessage {
	return &SMessage{
		FIncoming:  pIsIncoming,
		FMessage:   pMessage,
		FTimestamp: uint64(time.Now().Unix()), //nolint:gosec
	}
}

func LoadMessage(pMsgBytes []byte) IMessage {
	msg := &SMessage{}
	if err := json.Unmarshal(pMsgBytes, msg); err != nil {
		return nil
	}
	return msg
}

func (p *SMessage) IsIncoming() bool {
	return p.FIncoming
}

func (p *SMessage) GetMessage() string {
	return p.FMessage
}

func (p *SMessage) GetTimestamp() string {
	return time.Unix(int64(p.FTimestamp), 0).Format("2006-01-02T15:04:05") //nolint:gosec
}

func (p *SMessage) ToBytes() []byte {
	return encoding.SerializeJSON(p)
}

func (p *SMessage) ToString() string {
	return string(p.ToBytes())
}
