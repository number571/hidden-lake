package db

import (
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func GetChat(user *models.User, comp string) *models.Chat {
	var (
		err  error
		msg  models.Message
		chat = new(models.Chat)
	)
	id := GetUserId(user.Auth.Hashpasw)
	if id < 0 {
		return nil
	}
	rows, err := settings.DB.Query(
		"SELECT Name, Message, LastTime FROM Chat WHERE IdUser=$1 AND Companion=$2",
		id,
		comp,
	)
	if err != nil {
		panic("query 'getchat' failed")
	}
	defer rows.Close()
	chat.Companion = comp
	for rows.Next() {
		err = rows.Scan(
			&msg.Name,
			&msg.Text,
			&msg.Time,
		)
		if err != nil {
			break
		}
		decryptMessage(user, &msg)
		chat.Messages = append(chat.Messages, msg)
	}
	return chat
}
