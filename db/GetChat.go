package db

import (
	"../models"
	"../settings"
)

func GetChat(user *models.User, comp string) *models.Chat {
	var (
		err  error
		msg  models.Message
		chat = new(models.Chat)
	)
	rows, err := settings.DB.Query(
		"SELECT Name, Text, Time FROM Chat WHERE Hashname=$1 AND Companion=$2",
		user.Hashname,
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
