package dto

import (
	"encoding/json"
	"time"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/hidden-lake/internal/utils/chars"
)

var (
	_ IMessageBody = &SMessageBody{}
)

type SMessageBody struct {
	FData      string `json:"d"`
	FTimestamp uint64 `json:"t"`
}

func NewMessageBody(pData string, pTime time.Time) IMessageBody {
	msgBody := &SMessageBody{
		FData:      pData,
		FTimestamp: uint64(pTime.Unix()), //nolint:gosec
	}
	if err := ValidateMessageBody(msgBody); err != nil {
		panic(err)
	}
	return msgBody
}

func ValidateMessageBody(pBody IMessageBody) error {
	if chars.HasNotGraphicCharacters(pBody.GetData()) {
		return ErrHasNotGraphicChars
	}
	return nil
}

func LoadMessageBody(pData interface{}) (IMessageBody, error) {
	var msgBodyBytes []byte

	switch x := pData.(type) {
	case []byte:
		msgBodyBytes = x
	case string:
		msgBodyBytes = []byte(x)
	default:
		return nil, ErrUnknownType
	}

	msgBody := &SMessageBody{}
	if err := json.Unmarshal(msgBodyBytes, msgBody); err != nil {
		return nil, err
	}

	if err := ValidateMessageBody(msgBody); err != nil {
		return nil, err
	}

	return msgBody, nil
}

func (p *SMessageBody) GetData() string {
	return p.FData
}

func (p *SMessageBody) GetTimestamp() string {
	return time.Unix(int64(p.FTimestamp), 0).Format(time.DateTime) //nolint:gosec
}

func (p *SMessageBody) ToBytes() []byte {
	return encoding.SerializeJSON(p)
}

func (p *SMessageBody) ToString() string {
	return string(p.ToBytes())
}
