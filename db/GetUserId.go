package db

import (
	"github.com/number571/hiddenlake/settings"
)

func GetUserId(hashpasw string) int {
	var id = -2
	row := settings.DB.QueryRow("SELECT Id FROM User WHERE Hashpasw=$1", hashpasw)
	err := row.Scan(&id)
	if err != nil {
		return -1
	}
	return id
}
