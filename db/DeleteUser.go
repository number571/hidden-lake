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
	_, err = settings.DB.Exec("DELETE FROM Chat WHERE Hashname=$1", user.Hashname)
	if err != nil {
		panic("exec 'deleteuser.chat' failed")
	}
	_, err = settings.DB.Exec("DELETE FROM Client WHERE Contributor=$1", user.Hashname)
	if err != nil {
		panic("exec 'deleteuser.chat' failed")
	}
	_, err = settings.DB.Exec("DELETE FROM File WHERE Owner=$1", user.Hashname)
	if err != nil {
		panic("exec 'deleteuser.chat' failed")
	}
	return nil
}
