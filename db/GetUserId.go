package db

import (
	"github.com/number571/hiddenlake/settings"
)

func GetUserId(username string) int {
	var id = -2
	row := settings.DB.QueryRow("SELECT Id FROM User WHERE Username=$1", username)
	err := row.Scan(&id)
	if err != nil {
		return -1
	}
	return id
}
