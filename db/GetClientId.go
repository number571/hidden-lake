package db

import (
	"github.com/number571/hiddenlake/settings"
)

func GetClientId(iduser int, hashname string) int {
	var id = -2
	row := settings.DB.QueryRow("SELECT Id FROM Client WHERE IdUser=$1 AND Hashname=$2", iduser, hashname)
	err := row.Scan(&id)
	if err != nil {
		return -1
	}
	return id
}
