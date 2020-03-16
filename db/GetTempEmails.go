package db

import (
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func GetTempEmails(user *models.User, hashname string) []models.EmailType {
	var (
		emails []models.EmailType
		senderhash string
		public string 
		receiver string
		session string
		message string 
		salt string 
		hash string 
		sign string
		nonce uint64
	)
	id := GetUserId(user.Username)
	if id < 0 {
		return nil
	}
	rows, err := settings.DB.Query(
		"SELECT SenderHash, Sender, Receiver, Session, Message, Salt, Hash, Sign, Nonce FROM Email WHERE IdUser=$1 AND Receiver=$2 AND Temporary=1",
		id,
		hashname,
	)
	if err != nil {
		panic("query 'gettempemails' failed")
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&senderhash, &public, &receiver, &session, &message, &salt, &hash, &sign, &nonce)
		if err != nil {
			break
		}
		emails = append(emails, models.EmailType{
			Head: models.EmailHead{
				Sender: models.EmailSender{
					Public: public,
					Hashname: senderhash,
				},
				Receiver: receiver,
				Session: session,
			},
			Body: models.EmailBody{
				Data: message,
				Desc: models.EmailDesc{
					Rand: salt,
					Hash: hash,
					Sign: sign,
					Nonce: nonce,
					Difficulty: settings.DIFFICULTY,
				},
			},
		})
	}
	return emails
}
