package scheme

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
)

func LoadParticipantKey(pSchemeType ISchemeType, pFriendKey string) layer2.IParticipantKey {
	switch pSchemeType {
	case CHybridScheme:
		return asymmetric.LoadPubKey(pFriendKey)
	case CSymmetricScheme:
		return symmetric.NewCipherGCM(encoding.HexDecode(pFriendKey))
	default:
		panic("unknown key type")
	}
}
