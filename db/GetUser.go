package db

import (
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
	"github.com/number571/hiddenlake/utils"
)

func GetUser(username, password string) *models.User {
	name := gopeer.Base64Encode(gopeer.HashSum([]byte(username)))
	if !InUsers(name) {
		return nil
	}
	var (
		salt, hpasw, key string
		err error
	)
	row := settings.DB.QueryRow("SELECT Salt, Hashpasw, PrivateKey FROM User WHERE Username=$1", name)
	err = row.Scan(&salt, &hpasw, &key)
	if err != nil {
		return nil
	}
	pasw := gopeer.HashSum([]byte(password + salt))
	hashpasw := gopeer.HashSum(pasw)
	for i := 1; i < (1 << 20); i++ {
		hashpasw = gopeer.HashSum(hashpasw)
	}
	base64hashpasw := gopeer.Base64Encode(hashpasw)
	if base64hashpasw != hpasw {
		return nil
	}
	priv := gopeer.ParsePrivate(string(gopeer.DecryptAES(pasw, gopeer.Base64Decode(key))))
	return &models.User{
		Username: name,
		Hashname: gopeer.HashPublic(&priv.PublicKey),
		Auth: models.Auth{
			Hashpasw: base64hashpasw,
			Pasw:     pasw,
			Salt:     salt,
		},
		Keys: models.Keys{
			Private: priv,
			Public:  &priv.PublicKey,
		},
		Session: models.Session{
			Time: utils.CurrentTime(),
		},
	}
}
