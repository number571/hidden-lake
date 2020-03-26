package db

import (
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func GetTempEmails(user *models.User, hashname string) []models.EmailType {
	var (
		emails    []models.EmailType
		senderpub string
		receiver  string
		session   string
		title     string
		message   string
		salt      string
		hash      string
		sign      string
		nonce     uint64
	)
	id := GetUserId(user.Username)
	if id < 0 {
		return nil
	}
	rows, err := settings.DB.Query(
		"SELECT SenderPub, Receiver, Session, Title, Message, Salt, Hash, Sign, Nonce FROM Email WHERE IdUser=$1 AND Receiver=$2 AND Temporary=1",
		id,
		hashname,
	)
	if err != nil {
		panic("query 'gettempemails' failed")
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&senderpub, &receiver, &session, &title, &message, &salt, &hash, &sign, &nonce)
		if err != nil {
			break
		}
		emails = append(emails, models.EmailType{
			Head: models.EmailHead{
				Sender: models.EmailSender{
					Public:   senderpub,
					Hashname: gopeer.HashPublic(gopeer.ParsePublic(senderpub)),
				},
				Receiver: receiver,
				Session:  session,
			},
			Body: models.EmailBody{
				Data: models.EmailData{
					Head: title,
					Body: message,
				},
				Desc: models.EmailDesc{
					Rand:       salt,
					Hash:       hash,
					Sign:       sign,
					Nonce:      nonce,
					Difficulty: settings.DIFFICULTY,
				},
			},
		})
	}
	return emails
}
