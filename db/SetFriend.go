package db

import (
	"errors"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func SetFriend(user *models.User, hashname string) error {
	id := GetUserId(user.Username)
	if id < 0 {
		return errors.New("User id undefined")
	}
	idFriend := GetFriendId(id, hashname)
	if idFriend >= 0 {
		return errors.New("Friend already exist")
	}
	_, err := settings.DB.Exec(
		"INSERT INTO Friend (IdUser, Hashname) VALUES ($1, $2)",
		id,
		hashname,
	)
	if err != nil {
		panic("exec 'setfriend' failed")
	}
	return nil
}
