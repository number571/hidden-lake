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
	FNotify    bool          `json:"n"`
	FState     IMessageState `json:"s"`
	FUID       string        `json:"u"`
	FBody      string        `json:"b"`
	FTimestamp uint64        `json:"t"`
}

func NewMessage(pState IMessageState, pBody string) IMessage {
	msg, _ := NewMessageWithParams(pState, pBody, NewMessageUID(), time.Now(), false)
	return msg
}

func NewMessageWithParams(pState IMessageState, pBody string, pUID string, pTime time.Time, pIsNotify bool) (IMessage, error) {
	if !MessageUIDIsValid(pUID) {
		return nil, ErrUIDSize
	}
	return &SMessage{
		FNotify:    pIsNotify,
		FUID:       pUID,
		FState:     pState,
		FBody:      pBody,
		FTimestamp: uint64(pTime.Unix()), //nolint:gosec
	}, nil
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

	if !MessageUIDIsValid(msg.GetUID()) {
		return nil, ErrUIDSize
	}

	state := msg.GetState()
	if state == cBeginFlagState || state >= cEndFlagState {
		return nil, ErrUnknownState
	}

	return msg, nil
}

func (p *SMessage) SetState(pState IMessageState) {
	p.FState = pState
}

func (p *SMessage) GetState() IMessageState {
	return p.FState
}

func (p *SMessage) GetUID() string {
	return p.FUID
}

func (p *SMessage) GetBody() string {
	return p.FBody
}

func (p *SMessage) GetTimestamp() string {
	return time.Unix(int64(p.FTimestamp), 0).Format(time.DateTime) //nolint:gosec
}

func (p *SMessage) GetIsNotify() bool {
	return p.FNotify
}

func (p *SMessage) ToBytes() []byte {
	return encoding.SerializeJSON(p)
}

func (p *SMessage) ToString() string {
	return string(p.ToBytes())
}
