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

	_, err := settings.DB.Exec(
		"DELETE FROM Friends WHERE IdUser=$1 AND Hashname=$2",
		id,
		hashname,
	)
	if err != nil {
		panic("exec 'setfriend.delete' failed")
	}

	_, err = settings.DB.Exec(
		"INSERT INTO Friends (IdUser, Hashname) VALUES ($1, $2)",
		id,
		hashname,
	)
	if err != nil {
		panic("exec 'setfriend' failed")
	}

	return nil
}
