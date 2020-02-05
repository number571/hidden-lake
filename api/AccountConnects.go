package api

import (
	"encoding/json"
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/db"
	"github.com/number571/hiddenlake/settings"
	"github.com/number571/hiddenlake/utils"
	"net/http"
	"strings"
)

func AccountConnects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data struct {
		State string `json:"state"`
	}

	switch r.Method {
	case "GET":
		accountConnectsGET(w, r)
	case "PATCH":
		accountConnectsPATCH(w, r)
	default:
		data.State = "Method should be GET"
		json.NewEncoder(w).Encode(data)
	}
}

type connect struct {
	Connected bool `json:"connected"`
	Address   string `json:"address"`
	Hashname  string `json:"hashname"`
	Public    string `json:"public_key"`
}

// List of all clients.
func accountConnectsGET(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State    string   `json:"state"`
		Connects []connect `json:"connects"`
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

	user := settings.Users[token]
	client, ok := settings.Listener.Clients[user.Hashname]
	if !ok {
		data.State = "Current client is not exist"
		json.NewEncoder(w).Encode(data)
		return
	}

	clients := db.GetAllClients(user)
	if client == nil {
		data.State = "Error load clients from database"
		json.NewEncoder(w).Encode(data)
		return
	}

	for _, c := range clients {
		if c.Hashname == user.Hashname {
			continue
		}
		data.Connects = append(data.Connects, connect{
			Connected: client.InConnections(c.Hashname),
			Address: c.Address,
			Hashname: c.Hashname,
			Public: gopeer.StringPublic(c.Public),
		})
	}

	json.NewEncoder(w).Encode(data)
}

// List of current connections.
func accountConnectsPATCH(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State    string   `json:"state"`
		Connects []string `json:"connects"`
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

	user := settings.Users[token]
	client, ok := settings.Listener.Clients[user.Hashname]
	if !ok {
		data.State = "Current client is not exist"
		json.NewEncoder(w).Encode(data)
		return
	}

	for hash := range client.Connections {
		data.Connects = append(data.Connects, hash)
	}

	json.NewEncoder(w).Encode(data)
}
