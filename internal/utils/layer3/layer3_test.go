package layer3

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/payload"
)

func TestLayer3(t *testing.T) {
	t.Parallel()

	netk := "network_key"
	body := []byte("hello, world!")
	msg := NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{
				FNetworkKey:   netk,
				FWorkSizeBits: 8,
			}),
		}),
		body,
	)

	rawMsg, err := layer1.LoadMessage(
		layer1.NewSettings(&layer1.SSettings{
			FWorkSizeBits: 8,
		}),
		msg.ToBytes(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	rawMsgBytes, err := ExtractMessageBody(rawMsg)
	if err != nil {
		t.Error(err)
		return
	}

	decMsg, err := layer1.LoadMessage(
		layer1.NewSettings(&layer1.SSettings{
			FNetworkKey: netk,
		}),
		rawMsgBytes,
	)
	if err != nil {
		t.Error(err)
		return
	}

	decMsgBytes, err := ExtractMessageBody(decMsg)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(decMsgBytes, body) {
		t.Error("invalid body bytes")
		return
	}
}

func TestInvalidLayer3(t *testing.T) {
	t.Parallel()

	msg1 := layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{}),
		}),
		payload.NewPayload32(1, []byte{123}),
	)
	if _, err := ExtractMessageBody(msg1); err == nil {
		t.Error("success extract message with invalid salt size")
		return
	}

	msg2 := layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{}),
		}),
		payload.NewPayload32(1, bytes.Join(
			[][]byte{random.NewRandom().GetBytes(cSaltSize), []byte{123}},
			[]byte{},
		)),
	)
	if _, err := ExtractMessageBody(msg2); err == nil {
		t.Error("success extract message with invalid crc32 sum")
		return
	}
}

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SAppError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}
