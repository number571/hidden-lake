package db

import (
	"errors"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func SetChat(user *models.User, chat *models.Chat) error {
	id := GetUserId(user.Username)
	if id < 0 {
		return errors.New("User id undefined")
	}

	for index := range chat.Messages {
		encryptMessage(user, &chat.Messages[index])
		_, err := settings.DB.Exec(
			"INSERT INTO Chat (IdUser, IdClient, Name, Message, LastTime) VALUES ($1, $2, $3, $4, $5)",
			id,
			GetClientId(id, chat.Companion),
			chat.Messages[index].Name,
			chat.Messages[index].Text,
			chat.Messages[index].Time,
		)
		if err != nil {
			panic("exec 'setuser' failed")
		}
	}
	return nil
}
