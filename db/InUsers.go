package db

import (
	"github.com/number571/hiddenlake/settings"
)

func InUsers(hashpasw string) bool {
	var (
		checkHash string
		err       error
	)
	row := settings.DB.QueryRow("SELECT Hashpasw FROM User WHERE Hashpasw=$1", hashpasw)
	err = row.Scan(&checkHash)
	if err != nil {
		return false
	}
	return len(checkHash) == settings.LEN_BASE64_SHA256
}
