package api

import (
	"encoding/json"
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/db"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
	"net/http"
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

// List of all clients.
func accountConnectsGET(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State    string    `json:"state"`
		Connects []models.Connect `json:"connects"`
	}

	var (
		token string
		client = new(gopeer.Client)
	)
	switch {
	case isTokenAuthError(w, r, &token): return
	case isLifeTokenError(w, r, token): return
	case isGetClientError(w, r, client, token): return
	}

	clients := db.GetAllClients(settings.Users[token])
	if client == nil {
		data.State = "Error load clients from database"
		json.NewEncoder(w).Encode(data)
		return
	}

	for _, cl := range clients {
		if cl.Hashname == client.Hashname {
			continue
		}

		pub1 := gopeer.StringPublic(cl.Public)
		pub2 := gopeer.StringPublic(cl.PublicRecv)
		if pub1 != pub2 {
			continue
		}

		data.Connects = append(data.Connects, models.Connect{
			Connected: client.InConnections(cl.Hashname),
			Address:   cl.Address,
			Hashname:  cl.Hashname,
			PublicKey: gopeer.StringPublic(cl.Public),
		})
	}

	json.NewEncoder(w).Encode(data)
}

// List of current connections.
func accountConnectsPATCH(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State    string   `json:"state"`
		Connects []models.Connect `json:"connects"`
	}

	var (
		token string
		client = new(gopeer.Client)
	)
	switch {
	case isTokenAuthError(w, r, &token): return
	case isLifeTokenError(w, r, token): return
	case isGetClientError(w, r, client, token): return
	}

	for hash, cl := range client.Connections {
		if hash == client.Hashname {
			continue
		}

		data.Connects = append(data.Connects, models.Connect{
			Connected: client.InConnections(hash),
			Hidden:    gopeer.HashPublic(cl.Public) != gopeer.HashPublic(cl.PublicRecv),
			Address:   cl.Address,
			Hashname:  hash,
			ThrowNode: gopeer.HashPublic(cl.Public),
			PublicKey: gopeer.StringPublic(cl.PublicRecv),
		})
	}

	json.NewEncoder(w).Encode(data)
}
