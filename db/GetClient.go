package db

import (
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func GetClient(user *models.User, hashname string) *models.Client {
	var (
		address string
		public  string
	)
	id := GetUserId(user.Auth.Hashpasw)
	if id < 0 {
		return nil
	}
	row := settings.DB.QueryRow(
		"SELECT Address, PublicKey FROM Client WHERE IdUser=$1 AND Hashname=$2",
		id,
		hashname,
	)
	err := row.Scan(&address, &public)
	if err != nil {
		return nil
	}
	return &models.Client{
		Hashname: hashname,
		Address:  string(gopeer.DecryptAES(user.Auth.Pasw, gopeer.Base64Decode(address))),
		Public:   gopeer.ParsePublic(public),
	}
}
