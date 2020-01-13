package db

import (
	"../models"
	"../settings"
)

func SetChat(user *models.User, chat *models.Chat) error {
	for index := range chat.Messages {
		encryptMessage(user, &chat.Messages[index])
		_, err := settings.DB.Exec(
			"INSERT INTO Chat (Hashname, Companion, Name, Text, Time) VALUES ($1, $2, $3, $4, $5)",
			user.Hashname,
			chat.Companion,
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
