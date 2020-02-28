package db

import (
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func DeleteUser(user *models.User) error {
	_, err := settings.DB.Exec("DELETE FROM User WHERE Username=$1", user.Username)
	if err != nil {
		panic("exec 'deleteuser.user' failed")
	}
	return nil
}
