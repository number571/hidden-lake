package api

import (
	// "fmt"
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

type netData struct {
	List []models.LastMessage `json:"list"`
	Chat *models.Chat         `json:"chat"`
}

func NetworkChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data struct {
		State string `json:"state"`
	}

	switch r.Method {
	case "GET":
		networkGET(w, r)
	case "POST":
		networkPOST(w, r)
	case "DELETE":
		networkDELETE(w, r)
	default:
		data.State = "Method should be GET or POST"
		json.NewEncoder(w).Encode(data)
	}
}

// Delete chat.
func networkDELETE(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State string `json:"state"`
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

	if user.Hashname == read.Hashname {
		data.State = "Can't delete own chat"
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

// Send message.
func networkPOST(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State string `json:"state"`
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

// Get chat.
func networkGET(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State   string  `json:"state"`
		NetData netData `json:"netdata"`
	}

	var read struct {
		Hashname string `json:"hashname"`
	}

	read.Hashname = strings.Replace(r.URL.Path, "/api/network/chat/", "", 1)

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

	user := settings.Users[token]
	data.NetData.List = db.GetLastMessages(user)

	switch read.Hashname {
	case "", "null", "undefined":
		data.NetData.Chat = new(models.Chat)
	default:
		data.NetData.Chat = db.GetChat(user, read.Hashname)
	}

	json.NewEncoder(w).Encode(data)
}
