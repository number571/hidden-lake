package db

import (
	"errors"
	"github.com/number571/hiddenlake/utils"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func DeleteChat(user *models.User, hashname string) error {
	id := GetUserId(user.Auth.Hashpasw)
	if id < 0 {
		return errors.New("User id undefined")
	}
	_, err := settings.DB.Exec(
		"DELETE FROM Chat WHERE IdUser=$1 AND Companion=$2",
		id,
		hashname,
	)
	if err != nil {
		panic("exec 'deleteuser.chat' failed")
	}
	var(
		message = "chat is cleared"
		time    = utils.CurrentTime()
	)
	SetChat(user, &models.Chat{
		Companion: hashname,
		Messages: []models.Message{
			models.Message{
				Name: hashname,
				Text: message,
				Time: time,
			},
		},
	})
	return nil
}
