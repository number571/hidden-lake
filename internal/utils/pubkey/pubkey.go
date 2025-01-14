package pubkey

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
)

// echo PubKey{...} | sha384sum
func GetPubKeyHash(pPubKey asymmetric.IPubKey) string {
	return hashing.NewHasher([]byte(pPubKey.ToString())).ToString()
}
