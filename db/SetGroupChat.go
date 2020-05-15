package db

import (
	"errors"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func SetGroupChat(user *models.User, chat *models.Chat) error {
	id := GetUserId(user.Username)
	if id < 0 {
		return errors.New("User id undefined")
	}
	for index := range chat.Messages {
		encryptMessage(user, &chat.Messages[index])
		_, err := settings.DB.Exec(
			"INSERT INTO GlobalChat (IdUser, Founder, Name, Message, LastTime) VALUES ($1, $2, $3, $4, $5)",
			id,
			chat.Companion,
			chat.Messages[index].Name,
			chat.Messages[index].Text,
			chat.Messages[index].Time,
		)
		if err != nil {
			panic("exec 'setglobalchat' failed")
		}
	}
	return nil
}
