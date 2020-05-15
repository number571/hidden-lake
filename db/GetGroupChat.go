package db

import (
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func GetGroupChat(user *models.User, hashname string) *models.Chat {
	var (
		err  error
		msg  models.Message
		chat = new(models.Chat)
	)
	id := GetUserId(user.Username)
	if id < 0 {
		return nil
	}
	rows, err := settings.DB.Query(
		"SELECT Name, Message, LastTime FROM GlobalChat WHERE IdUser=$1 AND Founder=$2",
		id,
		hashname,
	)
	if err != nil {
		panic("query 'getglobalchat' failed")
	}
	defer rows.Close()
	chat.Companion = hashname
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
