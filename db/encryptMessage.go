package db

import (
	"../models"
	"github.com/number571/gopeer"
)

func encryptMessage(user *models.User, msg *models.Message) {
	msg.Name = gopeer.Base64Encode(
		gopeer.EncryptAES(
			user.Auth.Pasw,
			[]byte(msg.Name),
		),
	)
	msg.Text = gopeer.Base64Encode(
		gopeer.EncryptAES(
			user.Auth.Pasw,
			[]byte(msg.Text),
		),
	)
	msg.Time = gopeer.Base64Encode(
		gopeer.EncryptAES(
			user.Auth.Pasw,
			[]byte(msg.Time),
		),
	)
}
