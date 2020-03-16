package db

import (
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func GetEmail(user *models.User, hash string) *models.Email {
	var (
		incoming bool 
		lasttime string 
		hashname string
		public string 
		receiver string 
		message string 
		salt string
		sign string
		nonce uint64
	)
	id := GetUserId(user.Username)
	if id < 0 {
		return nil
	}
	row := settings.DB.QueryRow(
		"SELECT Incoming, LastTime, SenderHash, Sender, Receiver, Message, Salt, Sign, Nonce FROM Email WHERE IdUser=$1 AND Hash=$2 AND Temporary=0",
		id,
		hash,
	)
	err := row.Scan(&incoming, &lasttime, &hashname, &public, &receiver, &message, &salt, &sign, &nonce)
	if err != nil {
		return nil
	}
	return &models.Email{
		Info: models.EmailInfo{
			Incoming: incoming,
			Time: lasttime,
		},
		Email: models.EmailType{
			Head: models.EmailHead{
				Sender: models.EmailSender{
					Public: public,
					Hashname: hashname,
				},
				Receiver: receiver,
			},
			Body: models.EmailBody{
				Data: string(gopeer.DecryptAES(user.Auth.Pasw, gopeer.Base64Decode(message))),
				Desc: models.EmailDesc{
					Rand: salt,
					Hash: hash,
					Sign: sign,
					Nonce: nonce,
					Difficulty: settings.DIFFICULTY,
				},
			},
		},
	}
}