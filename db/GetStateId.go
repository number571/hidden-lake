package db

import (
	"github.com/number571/hiddenlake/settings"
)

func GetStateId(iduser int) int {
	var id = -2
	row := settings.DB.QueryRow("SELECT Id FROM State WHERE IdUser=$1", iduser)
	err := row.Scan(&id)
	if err != nil {
		return -1
	}
	return id
}
