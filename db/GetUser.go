package db

import (
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
	"github.com/number571/hiddenlake/utils"
)

func GetUser(pasw []byte) *models.User {
	hashpasw := gopeer.Base64Encode(gopeer.HashSum(pasw))
	if !InUsers(hashpasw) {
		return nil
	}
	var (
		key string
		err error
	)
	row := settings.DB.QueryRow("SELECT PrivateKey FROM User WHERE Hashpasw=$1", hashpasw)
	err = row.Scan(&key)
	if err != nil {
		return nil
	}
	priv := gopeer.ParsePrivate(string(gopeer.DecryptAES(pasw, gopeer.Base64Decode(key))))
	return &models.User{
		Hashname: gopeer.HashPublic(&priv.PublicKey),
		Auth: models.Auth{
			Hashpasw: hashpasw,
			Pasw:     pasw,
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
