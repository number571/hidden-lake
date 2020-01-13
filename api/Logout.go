package api

import (
	"../settings"
	"encoding/json"
	"net/http"
	"strings"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data struct {
		State string `json:"state"`
	}

	if r.Method != "POST" {
		data.State = "Method should be POST"
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

	hash := settings.Users[token].Hashname
	delete(settings.Listener.Clients, hash)
	delete(settings.Tokens, hash)
	delete(settings.Users, token)

	json.NewEncoder(w).Encode(data)
}
