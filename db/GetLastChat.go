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
	row := settings.DB.QueryRow(
		"SELECT Name, Text, Time FROM Chat WHERE Hashname=$1 AND Companion=$2 DESC",
		user.Hashname,
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
