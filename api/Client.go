package api

import (
	"../db"
	"../settings"
	"../utils"
	"../models"
	"encoding/json"
	"github.com/number571/gopeer"
	"net/http"
	"strings"
)

func Client(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data struct {
		State     string `json:"state"`
	}

	switch r.Method {
	case "GET":
		clientGET(w, r)
		return
	case "POST":
		clientPOST(w, r)
		return
	case "DELETE":
		clientDELETE(w, r)
		return
	}

	data.State = "Method should be GET"
	json.NewEncoder(w).Encode(data)
}

// Disconnect from client.
func clientDELETE(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State string `json:"state"`
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

// Connect to another client.
func clientPOST(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State string `json:"state"`
	}

	var read struct {
		Address   string `json:"address"`
		PublicKey string `json:"public_key"`
	}

	err := json.NewDecoder(r.Body).Decode(&read)
	if err != nil {
		data.State = "Error decode json format"
		json.NewEncoder(w).Encode(data)
		return
	}

	if len(strings.Split(read.Address, ":")) != 2 {
		data.State = "Address is not corrected"
		json.NewEncoder(w).Encode(data)
		return
	}

	public := gopeer.ParsePublic(read.PublicKey)
	if public == nil {
		data.State = "Error decode public key"
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
	client, ok := settings.Listener.Clients[user.Hashname]
	if !ok {
		data.State = "Current client is not exist"
		json.NewEncoder(w).Encode(data)
		return
	}

	dest := &gopeer.Destination{
		Address: read.Address,
		Public:  public,
	}
	err = client.Connect(dest)
	if err != nil {
		data.State = "Connect error"
		json.NewEncoder(w).Encode(data)
		return
	}

	hash := gopeer.HashPublic(public)
	err = db.SetClient(user, &models.Client{
		Hashname: hash,
		Address:  read.Address,
		Public:   public,
	})
	if err != nil {
		data.State = "Set client error"
		json.NewEncoder(w).Encode(data)
		return
	}

	message := "connection created"
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

	err = db.SetChat(user, &models.Chat{
		Companion: hash,
		Messages: []models.Message{
			models.Message{
				Name: hash,
				Text: message,
				Time: utils.CurrentTime(),
			},
		},
	})
	if err != nil {
		data.State = "Set chat error"
		json.NewEncoder(w).Encode(data)
		return
	}

	json.NewEncoder(w).Encode(data)
}

// Get client public information.
func clientGET(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Connected bool   `json:"connected"`
		Address   string `json:"address"`
		Hashname  string `json:"hashname"`
		PublicKey string `json:"public_key"`
		State     string `json:"state"`
	}

	var read struct {
		Hashname string `json:"hashname"`
	}

	read.Hashname = strings.Replace(r.URL.Path, "/api/network/client/", "", 1)

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