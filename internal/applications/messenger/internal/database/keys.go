package database

import (
	"fmt"

	"github.com/number571/go-peer/pkg/crypto/hashing"
)

const (
	cKeySizeTemplate          = "database[%s].friends[%s].size"
	cKeyMessageByEnumTemplate = "database[%s].friends[%s].messages[enum=%d]"
)

func getKeySize(pR IRelation) []byte {
	return []byte(fmt.Sprintf(
		cKeySizeTemplate,
		hashing.NewHasher([]byte(pR.IAm().ToString())).ToString(),
		hashing.NewHasher([]byte(pR.Friend().ToString())).ToString(),
	))
}

func getKeyMessageByEnum(pR IRelation, pI uint64) []byte {
	return []byte(fmt.Sprintf(
		cKeyMessageByEnumTemplate,
		hashing.NewHasher([]byte(pR.IAm().ToString())).ToString(),
		hashing.NewHasher([]byte(pR.Friend().ToString())).ToString(),
		pI,
	))
}
