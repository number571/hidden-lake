package api

import (
	"../settings"
	"../utils"
	"encoding/json"
	"github.com/number571/gopeer"
	"net/http"
	"strings"
)

func Account(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data struct {
		Address   string `json:"address"`
		Hashname  string `json:"hashname"`
		PublicKey string `json:"public_key"`
		State     string `json:"state"`
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
	err := settings.CheckLifetimeToken(token)
	if err != nil {
		data.State = "Token lifetime is over"
		json.NewEncoder(w).Encode(data)
		return
	} else {
		settings.Users[token].Session.Time = utils.CurrentTime()
	}

	hash := settings.Users[token].Hashname
	client, ok := settings.Listener.Clients[hash]
	if !ok {
		data.State = "Current client is not exist"
		json.NewEncoder(w).Encode(data)
		return
	}

	data.Address = client.Address
	data.Hashname = hash
	data.PublicKey = gopeer.StringPublic(settings.Users[token].Keys.Public)

	json.NewEncoder(w).Encode(data)
}
