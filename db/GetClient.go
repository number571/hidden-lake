package db

import (
	"../models"
	"../settings"
	"github.com/number571/gopeer"
)

func GetClient(user *models.User, hashname string) *models.Client {
	var (
		address string
		public  string
	)
	row := settings.DB.QueryRow(
		"SELECT Address, Public FROM Client WHERE Contributor=$1 AND Hashname=$2",
		user.Hashname,
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
