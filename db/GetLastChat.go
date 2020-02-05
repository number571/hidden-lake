package db

import (
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func GetLastChat(comp string, user *models.User) *models.Chat {
	var (
		msg  models.Message
		chat = new(models.Chat)
	)
	id := GetUserId(user.Auth.Hashpasw)
	if id < 0 {
		return nil
	}
	row := settings.DB.QueryRow(
		"SELECT Name, Message, LastTime FROM Chat WHERE IdUser=$1 AND Companion=$2 DESC",
		id,
		comp,
	)
	chat.Companion = comp
	err := row.Scan(
		&msg.Name,
		&msg.Text,
		&msg.Time,
	)
	if err != nil {
		return nil
	}
	decryptMessage(user, &msg)
	chat.Messages = append(chat.Messages, msg)
	return chat
}
