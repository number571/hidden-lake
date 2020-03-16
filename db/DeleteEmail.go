package db

import (
	"errors"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func DeleteEmail(user *models.User, hash string) error {
	id := GetUserId(user.Username)
	if id < 0 {
		return errors.New("User id undefined")
	}
	_, err := settings.DB.Exec("DELETE FROM Email WHERE IdUser=$1 AND Hash=$2",
		id,
		hash,
	)
	if err != nil {
		panic("exec 'deleteemail' failed")
	}
	return nil
}
