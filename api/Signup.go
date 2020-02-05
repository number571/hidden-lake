package api

import (
	"crypto/rsa"
	"encoding/json"
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/db"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
	"net/http"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data struct {
		State string `json:"state"`
	}

	if r.Method != "POST" {
		data.State = "Method should be POST"
		json.NewEncoder(w).Encode(data)
		return
	}

	var read struct {
		Username       string `json:"username"`
		Password       string `json:"password"`
		PasswordRepeat string `json:"password_repeat"`
		PrivateKey     string `json:"private_key"`
	}

	err := json.NewDecoder(r.Body).Decode(&read)
	if err != nil {
		data.State = "Error decode json format"
		json.NewEncoder(w).Encode(data)
		return
	}

	if read.Password != read.PasswordRepeat {
		data.State = "Passwords not equal"
		json.NewEncoder(w).Encode(data)
		return
	}

	user_len := len(read.Password)
	pasw_len := len(read.Password)
	if pasw_len < 6 || pasw_len > 128 ||
		user_len < 6 || user_len > 64 {
		data.State = "Length username or password does not match"
		json.NewEncoder(w).Encode(data)
		return
	}

	pasw := gopeer.HashSum([]byte(read.Username + read.Password))
	user := newUser(pasw, read.PrivateKey)
	if user == nil {
		data.State = "Error decode private key"
		json.NewEncoder(w).Encode(data)
		return
	}

	err = db.SetUser(user)
	if err != nil {
		data.State = "User already exist"
		json.NewEncoder(w).Encode(data)
		return
	}

	err = db.SetClient(user, &models.Client{
		Hashname: user.Hashname,
		Address:  settings.CFG.Host.Tcp.Ipv4 + settings.CFG.Host.Tcp.Port,
		Public:   user.Keys.Public,
	})
	if err != nil {
		data.State = "Set client error"
		json.NewEncoder(w).Encode(data)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func newUser(pasw []byte, private string) *models.User {
	var key *rsa.PrivateKey
	if private == "" {
		key = gopeer.GeneratePrivate(2048)
	} else {
		key = gopeer.ParsePrivate(private)
		if key == nil {
			return nil
		}
	}
	return &models.User{
		Hashname: gopeer.HashPublic(&key.PublicKey),
		Auth: models.Auth{
			Hashpasw: gopeer.Base64Encode(gopeer.HashSum(pasw)),
			Pasw:     pasw,
		},
		Keys: models.Keys{
			Private: key,
			Public:  &key.PublicKey,
		},
	}
}
