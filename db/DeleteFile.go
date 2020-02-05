package db

import (
	"errors"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
	"os"
)

func DeleteFile(user *models.User, filehash string) error {
	file := GetFile(user, filehash)
	if file == nil {
		return errors.New("File undefined")
	}

	id := GetUserId(user.Auth.Hashpasw)
	if id < 0 {
		return errors.New("User id undefined")
	}

	_, err := settings.DB.Exec("DELETE FROM File WHERE IdUser=$1 AND Hash=$2",
		id,
		filehash,
	)
	if err != nil {
		panic("exec 'deletefile' failed")
	}

	return os.Remove(settings.PATH_ARCHIVE + file.Path)
}
