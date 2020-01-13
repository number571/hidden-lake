package api

import (
	"../db"
	"../models"
	"../settings"
	"../utils"
	"encoding/json"
	"github.com/number571/gopeer"
	"net/http"
	"strings"
)

func Disconnect(w http.ResponseWriter, r *http.Request) {
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

	hash := settings.Users[token].Hashname
	client, ok := settings.Listener.Clients[hash]
	if !ok {
		data.State = "Current client is not exist"
		json.NewEncoder(w).Encode(data)
		return
	}

	if !client.InConnections(read.Hashname) {
		data.State = "User is not connected"
		json.NewEncoder(w).Encode(data)
		return
	}

	dest := &gopeer.Destination{
		Address: client.Connections[read.Hashname].Address,
		Public:  client.Connections[read.Hashname].Public,
	}

	message := "connection closed"
	_, err = client.SendTo(dest, &gopeer.Package{
		Head: gopeer.Head{
			Title:  settings.TITLE_MESSAGE,
			Option: settings.OPTION_GET,
		},
		Body: gopeer.Body{
			Data: message,
		},
	})
	if err != nil {
		data.State = "User can't receive message"
		json.NewEncoder(w).Encode(data)
		return
	}

	db.SetChat(settings.Users[token], &models.Chat{
		Companion: read.Hashname,
		Messages: []models.Message{
			models.Message{
				Name: hash,
				Text: message,
				Time: utils.CurrentTime(),
			},
		},
	})
	client.Disconnect(dest)

	json.NewEncoder(w).Encode(data)
}
