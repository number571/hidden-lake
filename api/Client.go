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

func Client(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data struct {
		Connected bool   `json:"connected"`
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

	var read struct {
		Hashname string `json:"hashname"`
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

	user := settings.Users[token]
	clientData := db.GetClient(user, read.Hashname)
	if clientData == nil {
		data.State = "Client undefined"
		json.NewEncoder(w).Encode(data)
		return
	}

	client, ok := settings.Listener.Clients[user.Hashname]
	if !ok {
		data.State = "Current client is not exist"
		json.NewEncoder(w).Encode(data)
		return
	}

	if client.InConnections(read.Hashname) {
		data.Connected = true
	}

	data.Address = clientData.Address
	data.Hashname = read.Hashname
	data.PublicKey = gopeer.StringPublic(clientData.Public)

	json.NewEncoder(w).Encode(data)
}
