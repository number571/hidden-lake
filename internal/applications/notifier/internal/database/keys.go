package database

import (
	"fmt"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

const (
	cKeyHashTemplate          = "database[%s].hash[%x]"
	cKeySizeTemplate          = "database[%s].channel[%s].size"
	cKeyMessageByEnumTemplate = "database[%s].channel[%s].messages[enum=%d]"
)

func getKeyHash(pPubKey asymmetric.IPubKey, pHash []byte) []byte {
	return []byte(fmt.Sprintf(
		cKeyHashTemplate,
		pPubKey.GetHasher().ToString(),
		pHash,
	))
}

func getKeySize(pR IRelation) []byte {
	return []byte(fmt.Sprintf(
		cKeySizeTemplate,
		pR.IAm().GetHasher().ToString(),
		pR.Key(),
	))
}

func getKeyMessageByEnum(pR IRelation, pI uint64) []byte {
	return []byte(fmt.Sprintf(
		cKeyMessageByEnumTemplate,
		pR.IAm().GetHasher().ToString(),
		pR.Key(),
		pI,
	))
}
