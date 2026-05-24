//go:build symmetric

package utils

import (
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
)

func LoadParticipantKey(pFriendKey string) layer2.IParticipantKey {
	return symmetric.NewCipherGCM(encoding.HexDecode(pFriendKey))
}
