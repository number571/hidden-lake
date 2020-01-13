package api

import (
	"../db"
	"../models"
	"../settings"
	"../utils"
	"encoding/json"
	"github.com/number571/gopeer"
	"golang.org/x/net/websocket"
	"net/http"
	"strings"
)

func Sendmsg(w http.ResponseWriter, r *http.Request) {
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
		Message  string `json:"message"`
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
	hash := user.Hashname
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

	message := strings.Replace(read.Message, "\n", " ", -1)
	_, err = client.Send(&gopeer.Package{
		To: gopeer.To{
			Receiver: gopeer.Receiver{
				Hashname: read.Hashname,
			},
			Address: client.Connections[read.Hashname].Address,
		},
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

	time := utils.CurrentTime()
	err = db.SetChat(user, &models.Chat{
		Companion: read.Hashname,
		Messages: []models.Message{
			models.Message{
				Name: hash,
				Text: message,
				Time: time,
			},
		},
	})
	if err != nil {
		data.State = "Set chat error"
		json.NewEncoder(w).Encode(data)
		return
	}

	var wsdata = struct {
		Comp struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"comp"`
		Text string `json:"text"`
		Time string `json:"time"`
	}{
		Comp: struct {
			From string `json:"from"`
			To   string `json:"to"`
		}{
			From: hash,
			To:   read.Hashname,
		},
		Text: message,
		Time: time,
	}
	if user.Session.Socket != nil {
		websocket.JSON.Send(user.Session.Socket, wsdata)
	}

	json.NewEncoder(w).Encode(data)
}
