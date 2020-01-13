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

func Delchat(w http.ResponseWriter, r *http.Request) {
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
		Hashname string `json:"hashname"`
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

	err = db.DeleteChat(user, read.Hashname)
	if err != nil {
		data.State = "Can't delete chat"
		json.NewEncoder(w).Encode(data)
		return
	}

	json.NewEncoder(w).Encode(data)
}
