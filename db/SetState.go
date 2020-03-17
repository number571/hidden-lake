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

	idState := GetStateId(id)
	if idState >= 0 {
		_, err := settings.DB.Exec(
			"UPDATE State SET UsedF2F=$1 WHERE IdUser=$2",
			state.UsedF2F,
			id,
		)
		if err != nil {
			panic("exec 'setstate.update' failed")
		}
		return nil
	}

	_, err := settings.DB.Exec(
		"INSERT INTO State (IdUser, UsedF2F) VALUES ($1, $2)",
		id,
		state.UsedF2F,
	)
	if err != nil {
		panic("exec 'setstate.insert' failed")
	}

	return nil
}
