package database

import (
	"fmt"
)

const (
	cKeySizeTemplate          = "database.friends[%s].size"
	cKeyMessageByEnumTemplate = "database.friends[%s].messages[enum=%d]"
	cKeyMessageByUidTemplate  = "database.friends[%s].messages[uid=%s]"
)

func getKeySize(pR string) []byte {
	return []byte(fmt.Sprintf(
		cKeySizeTemplate,
		pR,
	))
}

func getKeyMessageByEnum(pR string, pE uint64) []byte {
	return []byte(fmt.Sprintf(
		cKeyMessageByEnumTemplate,
		pR,
		pE,
	))
}

func getKeyMessageByUID(pR string, pU string) []byte {
	return []byte(fmt.Sprintf(
		cKeyMessageByUidTemplate,
		pR,
		pU,
	))
}
