package db

import (
	"../models"
	"../settings"
	"github.com/number571/gopeer"
)

func GetAllFiles(user *models.User) []models.File {
	var (
		files []models.File
		hash  string
		name  string
		size  uint64
	)
	rows, err := settings.DB.Query(
		"SELECT Hash, Name, Size FROM File WHERE Owner=$1",
		user.Hashname,
	)
	if err != nil {
		panic("query 'getallfiles' failed")
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&hash, &name, &size)
		if err != nil {
			break
		}
		files = append(files, models.File{
			Hash: hash,
			Name: string(gopeer.DecryptAES(user.Auth.Pasw, gopeer.Base64Decode(name))),
			Size: size,
		})
	}
	return files
}
