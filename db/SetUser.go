package db

import (
	"errors"
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func SetUser(user *models.User) error {
	if GetUserId(user.Username) >= 0 {
		return errors.New("User already exist")
	}
	_, err := settings.DB.Exec(
		"INSERT INTO User (Username, Salt, Hashpasw, PrivateKey) VALUES ($1, $2, $3, $4)",
		user.Username,
		user.Auth.Salt,
		user.Auth.Hashpasw,
		gopeer.Base64Encode(
			gopeer.EncryptAES(
				user.Auth.Pasw,
				[]byte(gopeer.StringPrivate(user.Keys.Private)),
			),
		),
	)
	SetState(user, &models.State{
		UsedF2F: false,
	})
	if err != nil {
		panic("exec 'setuser' failed")
	}
	return nil
}
