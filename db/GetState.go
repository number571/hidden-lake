package db

import (
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func GetState(user *models.User) *models.State {
	var (
		f2f bool
	)
	id := GetUserId(user.Username)
	if id < 0 {
		return nil
	}
	row := settings.DB.QueryRow(
		"SELECT UsedF2F FROM State WHERE IdUser=$1",
		id,
	)
	err := row.Scan(&f2f)
	if err != nil {
		return nil
	}
	return &models.State{
		UsedF2F: f2f,
	}
}
