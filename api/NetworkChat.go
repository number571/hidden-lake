package api

import (
	"encoding/json"
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/db"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
	"github.com/number571/hiddenlake/utils"
	"golang.org/x/net/websocket"
	"net/http"
	"strings"
)

type netdata struct {
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
		data.State = "Method should be GET, POST or DELETE"
		json.NewEncoder(w).Encode(data)
	}
}

// Get chat.
func networkGET(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State   string  `json:"state"`
		NetData netdata `json:"netdata"`
	}

	var read struct {
		Hashname string `json:"hashname"`
	}

	read.Hashname = strings.Replace(r.URL.Path, "/api/network/chat/", "", 1)

	var token string
	switch {
	case isTokenAuthError(w, r, &token): return
	case isLifeTokenError(w, r, token): return
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

// Send message.
func networkPOST(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State string `json:"state"`
	}

	var read struct {
		Hashname string `json:"hashname"`
		Message  string `json:"message"`
	}

	var (
		client = new(gopeer.Client)
		token string
	)
	switch {
	case isTokenAuthError(w, r, &token): return
	case isLifeTokenError(w, r, token): return
	case isDecodeError(w, r, &read): return
	case isGetClientError(w, r, client, token): return
	case isNotInConnectionsError(w, r, client, read.Hashname): return
	}

	message := strings.Replace(read.Message, "\n", " ", -1)
	_, err := client.SendTo(client.Destination(read.Hashname), &gopeer.Package{
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

	user := settings.Users[token]
	time := utils.CurrentTime()
	err = db.SetChat(user, &models.Chat{
		Companion: read.Hashname,
		Messages: []models.Message{
			models.Message{
				Name: client.Hashname(),
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
			From: client.Hashname(),
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

// Delete chat.
func networkDELETE(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State string `json:"state"`
	}

	var (
		read = new(userdata)
		user = new(models.User)
		token string
	)
	switch {
	case isTokenAuthError(w, r, &token): return
	case isLifeTokenError(w, r, token): return
	case isDecodeError(w, r, read): return
	case isGetUserError(w, r, user, read): return
	}

	err := db.DeleteChat(user, read.Hashname)
	if err != nil {
		data.State = "Can't delete chat"
		json.NewEncoder(w).Encode(data)
		return
	}

	json.NewEncoder(w).Encode(data)
}
