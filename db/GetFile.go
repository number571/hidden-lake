package db

import (
	"../models"
	"../settings"
	"github.com/number571/gopeer"
)

func GetFile(user *models.User, filehash string) *models.File {
	var (
		name string
		path string
		size uint64
	)
	row := settings.DB.QueryRow(
		"SELECT Name, Path, Size FROM File WHERE Owner=$1 AND Hash=$2",
		user.Hashname,
		filehash,
	)
	err := row.Scan(&name, &path, &size)
	if err != nil {
		return nil
	}
	return &models.File{
		Hash: filehash,
		Name: string(gopeer.DecryptAES(user.Auth.Pasw, gopeer.Base64Decode(name))),
		Path: string(gopeer.DecryptAES(user.Auth.Pasw, gopeer.Base64Decode(path))),
		Size: size,
	}
}
