package db

import (
	"github.com/number571/hiddenlake/settings"
)

func GetClientHashname(id int) string {
	var hashname string
	row := settings.DB.QueryRow("SELECT Hashname FROM Client WHERE Id=$1", id)
	err := row.Scan(&hashname)
	if err != nil {
		return ""
	}
	return hashname
}
