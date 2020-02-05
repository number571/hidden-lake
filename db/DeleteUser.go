package db

import (
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func DeleteUser(user *models.User) error {
	_, err := settings.DB.Exec("DELETE FROM User WHERE Hashpasw=$1", user.Auth.Hashpasw)
	if err != nil {
		panic("exec 'deleteuser.user' failed")
	}
	return nil
}
