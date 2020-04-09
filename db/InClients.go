package db

import (
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func InClients(user *models.User, hashname string) bool {
	id := GetUserId(user.Username)
	if id < 0 {
		return false
	}
	var (
		address string
		err     error
	)
	row := settings.DB.QueryRow(
		"SELECT Address FROM Client WHERE IdUser=$1 AND Hashname=$2",
		id,
		hashname,
	)
	err = row.Scan(&address)
	if err != nil {
		return false
	}
	return true
}
