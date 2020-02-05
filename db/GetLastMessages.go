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
	id := GetUserId(user.Auth.Hashpasw)
	if id < 0 {
		return nil
	}
	rows, err := settings.DB.Query(`
SELECT Companion, Name, Message, LastTime FROM (
    SELECT * FROM Chat WHERE IdUser=$1 ORDER BY Id DESC
) GROUP BY Companion ORDER BY Id DESC
`,
		id,
	)
	if err != nil {
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
