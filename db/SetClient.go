package db

import (
	"errors"
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func SetClient(user *models.User, client *models.Client) error {
	if gopeer.HashPublic(client.Public) != client.Hashname {
		return errors.New("hashname is not derived from the public key")
	}
	_, err := settings.DB.Exec(
		"DELETE FROM Client WHERE Contributor=$1 AND Hashname=$2",
		user.Hashname,
		client.Hashname,
	)
	if err != nil {
		panic("exec 'setclient.delete' failed")
	}
	_, err = settings.DB.Exec(
		"INSERT INTO Client (Contributor, Hashname, Address, Public) VALUES ($1, $2, $3, $4)",
		user.Hashname,
		client.Hashname,
		gopeer.Base64Encode(
			gopeer.EncryptAES(
				user.Auth.Pasw,
				[]byte(client.Address),
			),
		),
		gopeer.StringPublic(client.Public),
	)
	if err != nil {
		panic("exec 'setclient.insert' failed")
	}
	return nil
}
