package db

import (
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func InClients(user *models.User, hashname string) bool {
	var (
		public string
		err    error
	)
	row := settings.DB.QueryRow(
		"SELECT Public FROM User WHERE Contributor=$1 AND Hashname=$2",
		user.Hashname,
		hashname,
	)
	err = row.Scan(&public)
	if err != nil {
		return false
	}
	return public != ""
}
