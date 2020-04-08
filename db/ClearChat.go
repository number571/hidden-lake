package db

import (
	"errors"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
	"github.com/number571/hiddenlake/utils"
)

func ClearChat(user *models.User, hashname string) error {
	id := GetUserId(user.Username)
	if id < 0 {
		return errors.New("User id undefined")
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
		panic("exec 'deletechat' failed")
	}
	SetChat(user, &models.Chat{
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
