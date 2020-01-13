package db

import (
	"../models"
	"github.com/number571/gopeer"
)

func decryptMessage(user *models.User, msg *models.Message) {
	msg.Name = string(gopeer.DecryptAES(user.Auth.Pasw, gopeer.Base64Decode(msg.Name)))
	msg.Text = string(gopeer.DecryptAES(user.Auth.Pasw, gopeer.Base64Decode(msg.Text)))
	msg.Time = string(gopeer.DecryptAES(user.Auth.Pasw, gopeer.Base64Decode(msg.Time)))
}
