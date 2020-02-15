package api

import (
	"encoding/json"
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/db"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
	"net/http"
)

func AccountFriends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data struct {
		State string `json:"state"`
	}

	switch r.Method {
	case "GET":
		accountFriendsGET(w, r)
	case "POST":
		accountFriendsPOST(w, r)
	case "PATCH":
		accountFriendsPATCH(w, r)
	case "DELETE":
		accountFriendsDELETE(w, r)
	default:
		data.State = "Method should be GET, POST, PATCH or DELETE"
		json.NewEncoder(w).Encode(data)
	}
}

// List of friends.
func accountFriendsGET(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State    string    `json:"state"`
		StateF2F bool      `json:"statef2f"`
		Friends  []string  `json:"friends"`
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

	user := settings.Users[token]
	data.StateF2F = user.UsedF2F
	data.Friends = client.GetFriends()
	json.NewEncoder(w).Encode(data)
}

// Turn on/off check friend list.
func accountFriendsPOST(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State    string `json:"state"`
		StateF2F bool   `json:"statef2f"`
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

	user := settings.Users[token]
	f2f := !user.UsedF2F
	
	err := db.SetState(user, &models.State{
		UsedF2F: f2f,
	})
	if err != nil {
		data.State = "Can't set state"
		json.NewEncoder(w).Encode(data)
		return
	}

	client = settings.Listener.Clients[user.Hashname]

	user.UsedF2F  = f2f
	client.SetFriends(user.UsedF2F, client.GetFriends()...)

	data.StateF2F = user.UsedF2F
	json.NewEncoder(w).Encode(data)
}

// Append hash to friends.
func accountFriendsPATCH(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State string `json:"state"`
	}

	var read struct {
		Hashname string `json:"hashname"`
	}

	var (
		token string
		client = new(gopeer.Client)
	)
	switch {
	case isTokenAuthError(w, r, &token): return
	case isLifeTokenError(w, r, token): return
	case isDecodeError(w, r, &read): return
	case isGetClientError(w, r, client, token): return
	}

	if len(read.Hashname) != len(client.Hashname) {
		data.State = "Hashname length /= len(hash(public_key))"
		json.NewEncoder(w).Encode(data)
		return
	}

	if read.Hashname == client.Hashname {
		data.State = "Can't set friend hashname of current user"
		json.NewEncoder(w).Encode(data)
		return
	}

	user := settings.Users[token]
	err := db.SetFriend(user, read.Hashname)
	if err != nil {
		data.State = "Can't set friend"
		json.NewEncoder(w).Encode(data)
		return
	}

	client = settings.Listener.Clients[user.Hashname]

	client.AppendFriends(read.Hashname)
	json.NewEncoder(w).Encode(data)
}

// Delete hash from list of friends.
func accountFriendsDELETE(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State string `json:"state"`
	}

	var read struct {
		Hashname string `json:"hashname"`
	}

	var (
		token string
		client = new(gopeer.Client)
	)
	switch {
	case isTokenAuthError(w, r, &token): return
	case isLifeTokenError(w, r, token): return
	case isDecodeError(w, r, &read): return
	case isGetClientError(w, r, client, token): return
	}

	user := settings.Users[token]
	err := db.DeleteFriend(user, read.Hashname)
	if err != nil {
		data.State = "Can't delete friend"
		json.NewEncoder(w).Encode(data)
		return
	}

	client = settings.Listener.Clients[user.Hashname]

	client.DeleteFriends(read.Hashname)
	json.NewEncoder(w).Encode(data)
}
