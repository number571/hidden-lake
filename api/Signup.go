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

	var read = new(userdata)
	switch {
	case isDecodeError(w, r, read): return
	}

	user := newUser(read.Username, read.Password, read.PrivateKey)
	if user == nil {
		data.State = "Error decode private key"
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
	if user_len < 6 || user_len > 128 || pasw_len < 8 || pasw_len > 1024 {
		data.State = "Length username or password does not match"
		json.NewEncoder(w).Encode(data)
		return
	}

	err := db.SetUser(user)
	if err != nil {
		data.State = "User already exist"
		json.NewEncoder(w).Encode(data)
		return
	}

	err = db.SetClient(user, &models.Client{
		Hashname: user.Hashname,
		Address:  settings.CFG.Tcp.Ipv4 + settings.CFG.Tcp.Port,
		Public:   user.Keys.Public,
	})
	if err != nil {
		data.State = "Set client error"
		json.NewEncoder(w).Encode(data)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func newUser(username, password, private string) *models.User {
	var key *rsa.PrivateKey
	if private == "" {
		key = gopeer.GeneratePrivate(gopeer.Get("KEY_SIZE").(uint16))
	} else {
		key = gopeer.ParsePrivate(private)
		if key == nil {
			return nil
		}
	}
	salt := gopeer.Base64Encode(gopeer.GenerateRandomBytes(16))
	pasw := gopeer.HashSum([]byte(password + salt))

	hashpasw := gopeer.HashSum(pasw)
	for i := 1; i < (1 << 20); i++ {
		hashpasw = gopeer.HashSum(hashpasw)
	}

	return &models.User{
		Hashname: gopeer.HashPublic(&key.PublicKey),
		Username: gopeer.Base64Encode(gopeer.HashSum([]byte(username))),
		Auth: models.Auth{
			Hashpasw: gopeer.Base64Encode(hashpasw),
			Pasw:     pasw,
			Salt:     salt,
		},
		Keys: models.Keys{
			Private: key,
			Public:  &key.PublicKey,
		},
	}
}
