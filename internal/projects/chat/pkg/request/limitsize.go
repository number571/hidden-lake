package request

import (
	"crypto/ed25519"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
)

var (
	nullSeed = []byte("00000000000000000000000000000000")
	chanKey  = asymmetric.NewPrivKey().GetPubKey()
	pubKey   = ed25519.NewKeyFromSeed(nullSeed)
)

func GetMessageLimitSize(pldSize uint64) uint64 {
	req := BuildRequest(chanKey, pubKey, "")
	reqLen := uint64(len(req.ToBytes()))
	if pldSize < (reqLen + encoding.CSizeUint64) {
		panic("payload limit < header size of message")
	}
	return pldSize - reqLen - encoding.CSizeUint64
}
