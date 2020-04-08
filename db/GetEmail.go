package db

import (
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func GetEmail(user *models.User, hash string) *models.Email {
	var (
		incoming  bool
		lasttime  string
		senderpub string
		receiver  string
		title     string
		message   string
		salt      string
		sign      string
		nonce     uint64
	)
	id := GetUserId(user.Username)
	if id < 0 {
		return nil
	}
	row := settings.DB.QueryRow(
		"SELECT Incoming, LastTime, SenderPub, Receiver, Title, Message, Salt, Sign, Nonce FROM Email WHERE IdUser=$1 AND Hash=$2 AND Temporary=0",
		id,
		hash,
	)
	err := row.Scan(&incoming, &lasttime, &senderpub, &receiver, &title, &message, &salt, &sign, &nonce)
	if err != nil {
		return nil
	}
	return &models.Email{
		Info: models.EmailInfo{
			Incoming: incoming,
			Time:     lasttime,
		},
		Email: models.EmailType{
			Head: models.EmailHead{
				Sender: models.EmailSender{
					Public:   senderpub,
					Hashname: gopeer.HashPublic(gopeer.ParsePublic(senderpub)),
				},
				Receiver: receiver,
			},
			Body: models.EmailBody{
				Data: models.EmailData{
					Title: string(gopeer.DecryptAES(user.Auth.Pasw, gopeer.Base64Decode(title))),
					Message: string(gopeer.DecryptAES(user.Auth.Pasw, gopeer.Base64Decode(message))),
				},
				Desc: models.EmailDesc{
					Rand:       string(gopeer.DecryptAES(user.Auth.Pasw, gopeer.Base64Decode(salt))),
					Hash:       hash,
					Sign:       sign,
					Nonce:      nonce,
					Difficulty: settings.DIFFICULTY,
				},
			},
		},
	}
}
