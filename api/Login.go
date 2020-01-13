package api

import (
	"../db"
	"../settings"
	"encoding/json"
	"github.com/number571/gopeer"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data struct {
		Token string `json:"token"`
		State string `json:"state"`
	}

	if r.Method != "POST" {
		data.State = "Method should be POST"
		json.NewEncoder(w).Encode(data)
		return
	}

	var read struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&read)
	if err != nil {
		data.State = "Error decode json format"
		json.NewEncoder(w).Encode(data)
		return
	}

	pasw := gopeer.HashSum([]byte(read.Username + read.Password))
	user := db.GetUser(pasw)
	if user == nil {
		data.State = "User undefined"
		json.NewEncoder(w).Encode(data)
		return
	}

	token := gopeer.Base64Encode(gopeer.GenerateRandomBytes(20))
	hash := user.Hashname

	settings.Users[token] = user
	settings.Tokens[hash] = token
	settings.Listener.NewClient(user.Keys.Private)

	data.Token = token

	json.NewEncoder(w).Encode(data)
}
