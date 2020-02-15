package api

import (
	"encoding/json"
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/db"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
	"net/http"
)

func Account(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data struct {
		State string `json:"state"`
	}

	switch r.Method {
	case "GET":
		accountGET(w, r)
	case "POST":
		accountPOST(w, r)
	case "DELETE":
		accountDELETE(w, r)
	default:
		data.State = "Method should be GET, POST or DELETE"
		json.NewEncoder(w).Encode(data)
	}
}

// Get public information.
func accountGET(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Address     string `json:"address"`
		Hashname    string `json:"hashname"`
		PublicKey   string `json:"public_key"`
		Certificate string `json:"certificate"`
		State       string `json:"state"`
	}

	var (
		client = new(gopeer.Client)
		token string
	)

	switch {
	case isTokenAuthError(w, r, &token): return
	case isLifeTokenError(w, r, token): return
	case isGetClientError(w, r, client, token): return
	}

	data.Address = client.Address
	data.Hashname = client.Hashname
	data.PublicKey = gopeer.StringPublic(client.Keys.Public)
	data.Certificate = string(settings.Listener.Certificate)

	json.NewEncoder(w).Encode(data)
}

// Get private key.
func accountPOST(w http.ResponseWriter, r *http.Request) {
	var data struct {
		PrivateKey string `json:"private_key"`
		State      string `json:"state"`
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

	data.PrivateKey = gopeer.StringPrivate(settings.Users[token].Keys.Private)
	json.NewEncoder(w).Encode(data)
}

// Delete account.
func accountDELETE(w http.ResponseWriter, r *http.Request) {
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

	err := db.DeleteUser(user)
	if err != nil {
		data.State = "User not deleted"
		json.NewEncoder(w).Encode(data)
		return
	}

	deleteUserAuth(user)
	json.NewEncoder(w).Encode(data)
}
