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
	clientid := GetClientId(id, hashname)
	if clientid < 0 {
		return errors.New("Client id undefined")
	}
	_, err := settings.DB.Exec(
		"DELETE FROM Chat WHERE IdUser=$1 AND IdClient=$2",
		id,
		clientid,
	)
	if err != nil {
		panic("exec 'deleteclient.chat' failed")
	}
	_, err = settings.DB.Exec(
		"DELETE FROM GlobalChat WHERE IdUser=$1 AND Founder=$2",
		id,
		hashname,
	)
	if err != nil {
		panic("exec 'deleteclient.globalchat' failed")
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
