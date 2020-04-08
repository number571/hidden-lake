package db

import (
	"errors"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
	"github.com/number571/hiddenlake/utils"
)

func ClearGlobalChat(user *models.User, hashname string) error {
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
		panic("exec 'clearglobalchat' failed")
	}
	SetGlobalChat(user, &models.Chat{
		Companion: hashname,
		Messages: []models.Message{
			models.Message{
				Name: hashname,
				Text: "chat is cleared",
				Time: utils.CurrentTime(),
			},
		},
	})
	return nil
}
