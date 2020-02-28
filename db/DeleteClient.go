package db

import (
	"errors"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func DeleteClient(user *models.User, hashname string) error {
	id := GetUserId(user.Username)
	if id < 0 {
		return errors.New("User id undefined")
	}
	if user.Hashname == hashname {
		return errors.New("Can't delete user from function delete client")
	}
	_, err := settings.DB.Exec(
		"DELETE FROM Chat WHERE IdUser=$1 AND Companion=$2",
		id,
		hashname,
	)
	if err != nil {
		panic("exec 'deleteclient.chat' failed")
	}
	_, err = settings.DB.Exec(
		"DELETE FROM Client WHERE IdUser=$1 AND Hashname=$2", 
		id,
		hashname,
	)
	if err != nil {
		panic("exec 'deleteclient.client' failed")
	}
	return nil
}
