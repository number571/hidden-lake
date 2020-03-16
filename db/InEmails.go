package db

import (
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func InEmails(user *models.User, hash string) bool {
	var (
		lasttime string
		err    error
	)
	id := GetUserId(user.Username)
	if id < 0 {
		return false
	}
	row := settings.DB.QueryRow(
		"SELECT LastTime FROM Email WHERE IdUser=$1 AND Hash=$2",
		id,
		hash,
	)
	err = row.Scan(&lasttime)
	if err != nil {
		return false
	}
	return true
}
