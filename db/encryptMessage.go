package db

import (
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/models"
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
