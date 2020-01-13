package api

import (
	"../db"
	"../settings"
	"../utils"
	"encoding/json"
	"github.com/number571/gopeer"
	"net/http"
	"strings"
)

func Deluser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data struct {
		PrivateKey string `json:"private_key"`
		State      string `json:"state"`
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

	token := r.Header.Get("Authorization")
	token = strings.Replace(token, "Bearer ", "", 1)
	if _, ok := settings.Users[token]; !ok {
		data.State = "Tokened user undefined"
		json.NewEncoder(w).Encode(data)
		return
	}
	err = settings.CheckLifetimeToken(token)
	if err != nil {
		data.State = "Token lifetime is over"
		json.NewEncoder(w).Encode(data)
		return
	} else {
		settings.Users[token].Session.Time = utils.CurrentTime()
	}

	pasw := gopeer.HashSum([]byte(read.Username + read.Password))
	user := db.GetUser(pasw)
	if user == nil {
		data.State = "User undefined"
		json.NewEncoder(w).Encode(data)
		return
	}

	if user.Auth.Hashpasw != settings.Users[token].Auth.Hashpasw {
		data.State = "Users not equal"
		json.NewEncoder(w).Encode(data)
		return
	}

	err = db.DeleteUser(user)
	if err != nil {
		data.State = "User not deleted"
		json.NewEncoder(w).Encode(data)
		return
	}

	hash := user.Hashname
	delete(settings.Listener.Clients, hash)
	delete(settings.Tokens, hash)
	delete(settings.Users, token)

	json.NewEncoder(w).Encode(data)
}
