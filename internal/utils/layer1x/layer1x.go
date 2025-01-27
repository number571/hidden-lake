package layer1x

import (
	"bytes"
	"hash/crc32"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	cHeadVal  = 0x01
	cSaltSize = 32
)

func NewMessage(pSett layer1.IConstructSettings, pBody []byte) layer1.IMessage {
	salt := random.NewRandom().GetBytes(cSaltSize)
	return layer1.NewMessage(
		pSett,
		payload.NewPayload32(
			crc32.Checksum(salt, crc32.IEEETable),
			bytes.Join([][]byte{salt, pBody}, []byte{}),
		),
	)
}

func ExtractMessage(pMsg layer1.IMessage) ([]byte, error) {
	pld := pMsg.GetPayload()
	body := pld.GetBody()
	if len(body) < cSaltSize {
		return nil, ErrInvalidBody
	}
	if pld.GetHead() != crc32.Checksum(body[:cSaltSize], crc32.IEEETable) {
		return nil, ErrInvalidHead
	}
	return body[cSaltSize:], nil
}
