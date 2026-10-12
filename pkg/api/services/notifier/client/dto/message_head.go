package dto

import (
	"encoding/json"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IMessageHead = &SMessageHead{}
)

type SMessageHead struct {
	FState IMessageState `json:"s"`
	FUID   string        `json:"u"`
}

func NewMessageHead(pState IMessageState, pUID string) IMessageHead {
	msgHead := &SMessageHead{
		FState: pState,
		FUID:   pUID,
	}
	if err := ValidateMessageHead(msgHead); err != nil {
		panic(err)
	}
	return msgHead
}

func LoadMessageHead(pData interface{}) (IMessageHead, error) {
	var msgHeadBytes []byte

	switch x := pData.(type) {
	case []byte:
		msgHeadBytes = x
	case string:
		msgHeadBytes = []byte(x)
	default:
		return nil, ErrUnknownType
	}

	msgHead := &SMessageHead{}
	if err := json.Unmarshal(msgHeadBytes, msgHead); err != nil {
		return nil, err
	}

	if err := ValidateMessageHead(msgHead); err != nil {
		return nil, err
	}

	return msgHead, nil
}

func ValidateMessageHead(pHead IMessageHead) error {
	if !MessageStateIsValid(pHead.GetState()) {
		return ErrUnknownState
	}
	if !MessageUIDIsValid(pHead.GetUID()) {
		return ErrUIDSize
	}
	return nil
}

func (p *SMessageHead) SetState(pState IMessageState) {
	p.FState = pState
}

func (p *SMessageHead) GetState() IMessageState {
	return p.FState
}

func (p *SMessageHead) GetUID() string {
	return p.FUID
}

func (p *SMessageHead) ToBytes() []byte {
	return encoding.SerializeJSON(p)
}

func (p *SMessageHead) ToString() string {
	return string(p.ToBytes())
}
