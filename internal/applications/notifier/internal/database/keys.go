package database

import (
	"fmt"
)

const (
	cKeySizeTemplate          = "database[%s].notifier.size"
	cKeyHashTemplate          = "database[%s].notifier.hash[%x]"
	cKeyMessageByEnumTemplate = "database[%s].notifier.messages[enum=%d]"
)

func getKeySize(pR IRelation) []byte {
	return []byte(fmt.Sprintf(
		cKeySizeTemplate,
		pR.IAm().GetHasher().ToString(),
	))
}

func getKeyHash(pR IRelation, pHash []byte) []byte {
	return []byte(fmt.Sprintf(
		cKeyHashTemplate,
		pR.IAm().GetHasher().ToString(),
		pHash,
	))
}

func getKeyMessageByEnum(pR IRelation, pI uint64) []byte {
	return []byte(fmt.Sprintf(
		cKeyMessageByEnumTemplate,
		pR.IAm().GetHasher().ToString(),
		pI,
	))
}
