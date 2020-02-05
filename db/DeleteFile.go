package db

import (
	"../models"
	"../settings"
	"errors"
	"os"
)

func DeleteFile(user *models.User, filehash string) error {
	file := GetFile(user, filehash)
	if file == nil {
		return errors.New("File undefined")
	}

	_, err := settings.DB.Exec("DELETE FROM File WHERE Owner=$1 AND Hash=$2",
		user.Hashname,
		filehash,
	)
	if err != nil {
		panic("exec 'deletefile' failed")
	}

	return os.Remove(settings.PATH_ARCHIVE + file.Path)
}
