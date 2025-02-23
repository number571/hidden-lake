package layer3

import (
	"bytes"
	"hash/crc32"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	cSaltSize = 32
)

func NewMessage(pSett layer1.IConstructSettings, pBody []byte) layer1.IMessage {
	cfgSett := pSett.GetSettings()
	encMsg := layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{
				FNetworkKey: cfgSett.GetNetworkKey(),
			}),
		}),
		newPayload(pBody),
	)
	return layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{
				FWorkSizeBits: cfgSett.GetWorkSizeBits(),
			}),
			FParallel: pSett.GetParallel(),
		}),
		newPayload(encMsg.ToBytes()),
	)
}

func ExtractMessageBody(pMsg layer1.IMessage) ([]byte, error) {
	pld := pMsg.GetPayload()
	body := pld.GetBody()
	if len(body) < cSaltSize {
		return nil, ErrInvalidBody
	}
	if pld.GetHead() != crc32.Checksum(body, crc32.IEEETable) {
		return nil, ErrInvalidHead
	}
	return body[cSaltSize:], nil
}

func newPayload(pBody []byte) payload.IPayload32 {
	saltBody := bytes.Join(
		[][]byte{random.NewRandom().GetBytes(cSaltSize), pBody},
		[]byte{},
	)
	return payload.NewPayload32(
		crc32.Checksum(saltBody, crc32.IEEETable),
		saltBody,
	)
}
