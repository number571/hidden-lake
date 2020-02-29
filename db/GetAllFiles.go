package db

import (
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func GetAllFiles(user *models.User) []models.File {
	var (
		files []models.File
		hash  string
		name  string
		size  uint64
		encr  bool
	)
	id := GetUserId(user.Username)
	if id < 0 {
		return nil
	}
	rows, err := settings.DB.Query(
		"SELECT Hash, Name, Size, Encr FROM File WHERE IdUser=$1",
		id,
	)
	if err != nil {
		panic("query 'getallfiles' failed")
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&hash, &name, &size, &encr)
		if err != nil {
			break
		}
		files = append(files, models.File{
			Hash: hash,
			Name: string(gopeer.DecryptAES(user.Auth.Pasw, gopeer.Base64Decode(name))),
			Size: size,
			Encr: encr,
		})
	}
	return files
}
