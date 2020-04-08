package handle

import (
	"bytes"
	"errors"
	"crypto/rsa"
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/db"
	"github.com/number571/hiddenlake/utils"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func NewEmail(client *gopeer.Client, public *rsa.PublicKey, title, message string) (*models.EmailType, string) {
	var (
		publicst = gopeer.StringPublic(client.Public())
		hashname = gopeer.HashPublic(gopeer.ParsePublic(publicst))
		receiver = gopeer.HashPublic(public)
		session  = gopeer.GenerateRandomBytes(32)
		random   = gopeer.GenerateRandomBytes(16)
		hash     = gopeer.HashSum(bytes.Join(
			[][]byte{
				[]byte(hashname),
				[]byte(receiver),
				[]byte(title),
				[]byte(message),
				random,
			},
			[]byte{},
		))
	)
	return &models.EmailType{
		Head: models.EmailHead{
			Sender: models.EmailSender{
				Public: publicst,
				Hashname: hashname,
			},
			Receiver: receiver,
			Session: gopeer.Base64Encode(gopeer.EncryptRSA(public, session)),
		},
		Body: models.EmailBody{
			Data: models.EmailData{
				Title: gopeer.Base64Encode(gopeer.EncryptAES(session, []byte(title))),
				Message: gopeer.Base64Encode(gopeer.EncryptAES(session, []byte(message))),
			},
			Desc: models.EmailDesc{
				Rand: gopeer.Base64Encode(gopeer.EncryptAES(session, random)),
				Hash: gopeer.Base64Encode(hash),
				Sign: gopeer.Base64Encode(gopeer.Sign(client.Private(), hash)),
				Nonce: gopeer.ProofOfWork(hash, settings.DIFFICULTY),
				Difficulty: settings.DIFFICULTY,
			},
		},
	}, gopeer.Base64Encode(random)
}

func ReadEmail(client *gopeer.Client, email *models.EmailType) (*models.Email, error) {
	session := gopeer.DecryptRSA(client.Private(), gopeer.Base64Decode(email.Head.Session))
	if session == nil {
		return nil, errors.New("error read session key")
	}
	title := string(gopeer.DecryptAES(session, gopeer.Base64Decode(email.Body.Data.Title)))
	if len(title) > 128 {
		return nil, errors.New("email.title size exceeded")
	}
	message := string(gopeer.DecryptAES(session, gopeer.Base64Decode(email.Body.Data.Message)))
	if len(message) >= settings.EMAIL_SIZE {
		return nil, errors.New("email.message size exceeded")
	}
	random := gopeer.DecryptAES(session, gopeer.Base64Decode(email.Body.Desc.Rand))
	hash := gopeer.HashSum(bytes.Join(
		[][]byte{
			[]byte(email.Head.Sender.Hashname),
			[]byte(email.Head.Receiver),
			[]byte(title),
			[]byte(message),
			random,
		},
		[]byte{},
	))
	if gopeer.Base64Encode(hash) != email.Body.Desc.Hash {
		return nil, errors.New("hashes not equal")
	}
	if email.Body.Desc.Difficulty != settings.DIFFICULTY {
		return nil, errors.New("difficulty does not match")
	}
	public := gopeer.ParsePublic(email.Head.Sender.Public)
	if public == nil {
		return nil, errors.New("error parse public")
	}
	if gopeer.HashPublic(public) != email.Head.Sender.Hashname {
		return nil, errors.New("hash(public) and hashname not equal")
	}
	if gopeer.Verify(public, hash, gopeer.Base64Decode(email.Body.Desc.Sign)) != nil {
		return nil, errors.New("email sign invalid")
	}
	if !gopeer.NonceIsValid(gopeer.Base64Decode(email.Body.Desc.Hash), uint(email.Body.Desc.Difficulty), email.Body.Desc.Nonce) {
		return nil, errors.New("email nonce is invalid")
	}
	return &models.Email{
		Info: models.EmailInfo{
			Incoming: true,
			Time: utils.CurrentTime(),
		},
		Email: models.EmailType{
			Head: models.EmailHead{
				Sender: email.Head.Sender,
				Receiver: email.Head.Receiver,
				Session: email.Head.Session,
			},
			Body: models.EmailBody{
				Data: models.EmailData{
					Title: title,
					Message: message,
				},
				Desc: models.EmailDesc{
					Rand: gopeer.Base64Encode(random),
					Hash: email.Body.Desc.Hash,
					Sign: email.Body.Desc.Sign,
					Nonce: email.Body.Desc.Nonce,
					Difficulty: email.Body.Desc.Difficulty,
				},
			},
		},
	}, nil
}

func getEmail(client *gopeer.Client, pack *gopeer.Package) (set string) {
	var (
		email = new(models.EmailType)
		token = settings.Tokens[client.Hashname()]
		user  = settings.Users[token]
	)
	gopeer.UnpackJSON([]byte(pack.Body.Data), email)
	if email == nil {
		return
	}
	// Send list of emails.
	if email.Head.Session == gopeer.Get("OPTION_GET").(string) {
		if !client.InConnections(email.Head.Receiver) {
			return
		}
		emails := db.GetTempEmails(user, email.Head.Receiver)
		return string(gopeer.PackJSON(emails))
	}
	// Check email's valid.
	if email.Head.Receiver == email.Head.Sender.Hashname {
		return
	}
	if email.Body.Desc.Difficulty != settings.DIFFICULTY {
		return
	}
	public := gopeer.ParsePublic(email.Head.Sender.Public)
	if public == nil {
		return
	}
	if gopeer.HashPublic(public) != email.Head.Sender.Hashname {
		return
	}
	if gopeer.Verify(public, gopeer.Base64Decode(email.Body.Desc.Hash), gopeer.Base64Decode(email.Body.Desc.Sign)) != nil {
		return
	}
	if !gopeer.NonceIsValid(gopeer.Base64Decode(email.Body.Desc.Hash), uint(email.Body.Desc.Difficulty), email.Body.Desc.Nonce) {
		return
	}
	if _, ok := client.F2F.Friends[email.Head.Sender.Hashname]; client.F2F.Perm && !ok {
		return
	}
	// Redirect or save email.
	switch {
	case client.Hashname() == email.Head.Receiver:
		newEmail, err := ReadEmail(client, email)
		if err != nil {
			return
		}
		db.SetEmail(user, models.IsPermEmail, newEmail)
	case client.InConnections(email.Head.Receiver):
		dest := client.Destination(email.Head.Receiver)
		client.SendTo(dest, &gopeer.Package{
			Head: gopeer.Head{
				Title: settings.TITLE_EMAIL,
				Option: gopeer.Get("OPTION_GET").(string),
			},
			Body: gopeer.Body{
				Data: pack.Body.Data,
			},
		})
	default:
		if _, ok := client.F2F.Friends[email.Head.Receiver]; client.F2F.Perm && !ok {
			return
		}
		db.SetEmail(user, models.IsTempEmail, &models.Email{
			Info: models.EmailInfo{
				Incoming: true,
				Time: utils.CurrentTime(),
			},
			Email: *email,
		})
	}
	return set
}

func setEmail(client *gopeer.Client, pack *gopeer.Package) {
	var (
		emails []models.EmailType
		token = settings.Tokens[client.Hashname()]
		user  = settings.Users[token]
	)
	gopeer.UnpackJSON([]byte(pack.Body.Data), &emails)
	if emails == nil {
		return
	}
	for _, email := range emails {
		if _, ok := client.F2F.Friends[email.Head.Sender.Hashname]; client.F2F.Perm && !ok {
			continue
		}
		newEmail, err := ReadEmail(client, &email)
		if err != nil {
			continue
		}
		db.SetEmail(user, models.IsPermEmail, newEmail)
	}
}
