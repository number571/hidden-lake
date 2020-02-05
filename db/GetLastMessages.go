package db

import (
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func GetLastMessages(user *models.User) []models.LastMessage {
	var (
		err  error
		msg  models.LastMessage
		msgs []models.LastMessage
	)
	rows, err := settings.DB.Query(`
SELECT Companion, Name, Text, Time FROM (
    SELECT * FROM Chat WHERE Hashname=$1 ORDER BY Id DESC
) GROUP BY Companion ORDER BY Id DESC
`,
		user.Hashname,
	)
	if err != nil {
		println(err)
		panic("query 'getlastmessages' failed")
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&msg.Companion,
			&msg.Message.Name,
			&msg.Message.Text,
			&msg.Message.Time,
		)
		if err != nil {
			break
		}
		decryptMessage(user, &msg.Message)
		msgs = append(msgs, msg)
	}
	return msgs
}
