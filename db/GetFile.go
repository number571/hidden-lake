package db

import (
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func GetFile(user *models.User, filehash string) *models.File {
	var (
		name string
		path string
		size uint64
	)
	id := GetUserId(user.Username)
	if id < 0 {
		return nil
	}
	row := settings.DB.QueryRow(
		"SELECT Name, Path, Size FROM File WHERE IdUser=$1 AND Hash=$2",
		id,
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
