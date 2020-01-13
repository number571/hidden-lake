package db

import (
	"../models"
	"../settings"
)

func DeleteChat(user *models.User, hashname string) error {
	_, err := settings.DB.Exec(
		"DELETE FROM Chat WHERE Hashname=$1 AND Companion=$2",
		user.Hashname,
		hashname,
	)
	if err != nil {
		panic("exec 'deleteuser.chat' failed")
	}
	_, err = settings.DB.Exec(
		"DELETE FROM Client WHERE Contributor=$1 AND Hashname=$2",
		user.Hashname,
		hashname,
	)
	if err != nil {
		panic("exec 'setclient.delete' failed")
	}
	return nil
}
