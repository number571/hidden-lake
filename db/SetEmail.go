package db

import (
	"errors"
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func SetEmail(user *models.User, option models.EmailSaveOption, email *models.Email) error {
	id := GetUserId(user.Username)
	if id < 0 {
		return errors.New("User id undefined")
	}
	if InEmails(user, email.Email.Body.Desc.Hash) {
		return errors.New("Email already exist")
	}
	var (
		message string
		session string
	)
	switch option {
	case models.IsTempEmail:
		session = email.Email.Head.Session
		message = email.Email.Body.Data
	case models.IsPermEmail:
		session = ""
		message = gopeer.Base64Encode(
			gopeer.EncryptAES(
				user.Auth.Pasw,
				[]byte(email.Email.Body.Data),
			),
		)
	}
	_, err := settings.DB.Exec(
		"INSERT INTO Email (IdUser, Incoming, Temporary, LastTime, SenderHash, Sender, Receiver, Session, Message, Salt, Hash, Sign, Nonce) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		id,
		email.Info.Incoming,
		option,
		email.Info.Time,
		email.Email.Head.Sender.Hashname,
		email.Email.Head.Sender.Public,
		email.Email.Head.Receiver,
		session,
		message,
		email.Email.Body.Desc.Rand,
		email.Email.Body.Desc.Hash,
		email.Email.Body.Desc.Sign,
		email.Email.Body.Desc.Nonce,
	)
	if err != nil {
		panic("exec 'setemail' failed")
	}
	return nil
}
