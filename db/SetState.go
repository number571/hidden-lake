package db

import (
	"errors"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func SetState(user *models.User, state *models.State) error {
	id := GetUserId(user.Username)
	if id < 0 {
		return errors.New("User id undefined")
	}

	_, err := settings.DB.Exec(
		"DELETE FROM State WHERE IdUser=$1",
		id,
	)
	if err != nil {
		panic("exec 'setstate.delete' failed")
	}

	_, err = settings.DB.Exec(
		"INSERT INTO State (IdUser, UsedF2F) VALUES ($1, $2)",
		id,
		state.UsedF2F,
	)
	if err != nil {
		panic("exec 'setstate' failed")
	}

	return nil
}
