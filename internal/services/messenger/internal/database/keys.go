package database

import (
	"fmt"
)

const (
	cKeySizeTemplate          = "database.friends[%s].size"
	cKeyMessageByEnumTemplate = "database.friends[%s].messages[enum=%d]"
)

func getKeySize(pR string) []byte {
	return []byte(fmt.Sprintf(
		cKeySizeTemplate,
		pR,
	))
}

func getKeyMessageByEnum(pR string, pI uint64) []byte {
	return []byte(fmt.Sprintf(
		cKeyMessageByEnumTemplate,
		pR,
		pI,
	))
}
