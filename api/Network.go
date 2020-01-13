package api

import (
	"../db"
	"../models"
	"../settings"
	"../utils"
	"encoding/json"
	"net/http"
	"strings"
)

type netData struct {
	List []models.LastMessage `json:"list"`
	Chat *models.Chat         `json:"chat"`
}

func Network(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data struct {
		State   string  `json:"state"`
		NetData netData `json:"netdata"`
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

	data.NetData.List = db.GetLastMessages(settings.Users[token])
	switch read.Hashname {
	case "":
		data.NetData.Chat = new(models.Chat)
	default:
		data.NetData.Chat = db.GetChat(read.Hashname, settings.Users[token])
	}

	json.NewEncoder(w).Encode(data)
}
