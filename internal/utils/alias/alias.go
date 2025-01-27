package alias

import (
	"bytes"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func GetAliasesList(pFriends map[string]asymmetric.IPubKey) []string {
	result := make([]string, 0, len(pFriends))
	for alias := range pFriends {
		result = append(result, alias)
	}
	return result
}

func GetAliasByPubKey(pFriends map[string]asymmetric.IPubKey, pPubKey asymmetric.IPubKey) string {
	for alias, pubKey := range pFriends {
		if bytes.Equal(pubKey.ToBytes(), pPubKey.ToBytes()) {
			return alias
		}
	}
	return ""
}
