package db

import (
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func SetFile(user *models.User, file *models.File) error {
	_, err := settings.DB.Exec(
		"INSERT INTO File (Owner, Hash, Name, Path, Size) VALUES ($1, $2, $3, $4, $5)",
		user.Hashname,
		file.Hash,
		gopeer.Base64Encode(
			gopeer.EncryptAES(
				user.Auth.Pasw,
				[]byte(file.Name),
			),
		),
		gopeer.Base64Encode(
			gopeer.EncryptAES(
				user.Auth.Pasw,
				[]byte(file.Path),
			),
		),
		file.Size,
	)
	if err != nil {
		panic("exec 'setfile' failed")
	}
	return nil
}
