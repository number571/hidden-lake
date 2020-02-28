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
	case "DELETE":
		accountConnectsDELETE(w, r)
	default:
		data.State = "Method should be GET, PATCH or DELETE"
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
	if clients == nil {
		data.State = "Error load clients from database"
		json.NewEncoder(w).Encode(data)
		return
	}

	for _, cl := range clients {
		if cl.Hashname == client.Hashname() {
			continue
		}

		if gopeer.HashPublic(cl.Public) != gopeer.HashPublic(cl.ThrowClient) {
			continue
		}

		data.Connects = append(data.Connects, models.Connect{
			Connected:   client.InConnections(cl.Hashname),
			Address:     cl.Address,
			Hashname:    cl.Hashname,
			Public:      gopeer.StringPublic(cl.Public),
			Certificate: cl.Certificate,
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
		if hash == client.Hashname() {
			continue
		}

		data.Connects = append(data.Connects, models.Connect{
			Connected:   client.InConnections(hash),
			Hidden:      gopeer.HashPublic(cl.Public()) != gopeer.HashPublic(cl.Throw()),
			Address:     cl.Address(),
			Hashname:    hash,
			Public:      gopeer.StringPublic(cl.Public()),
			ThrowClient: gopeer.HashPublic(cl.Throw()),
			Certificate: string(cl.Certificate()),
		})
	}

	json.NewEncoder(w).Encode(data)
}

// Delete client from user data.
func accountConnectsDELETE(w http.ResponseWriter, r *http.Request) {
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

	err := db.DeleteClient(user, read.Hashname)
	if err != nil {
		data.State = "Can't delete client"
		json.NewEncoder(w).Encode(data)
		return
	}

	json.NewEncoder(w).Encode(data)
}
