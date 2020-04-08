package db

import (
	"errors"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func DeleteGlobalChat(user *models.User, hashname string) error {
	id := GetUserId(user.Username)
	if id < 0 {
		return errors.New("User id undefined")
	}
	if user.Hashname == hashname {
		return errors.New("Delete own chat error")
	}
	_, err := settings.DB.Exec(
		"DELETE FROM GlobalChat WHERE IdUser=$1 AND Founder=$2",
		id,
		hashname,
	)
	if err != nil {
		panic("exec 'deleteglobalchat' failed")
	}

	delete(user.Temp.ChatMap.Member, hashname)
	return nil
}
