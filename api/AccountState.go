package api

import (
	"strings"
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/settings"
	"github.com/number571/hiddenlake/db"
	"encoding/json"
	"net/http"
)

func AccountState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data struct {
		State string `json:"state"`
	}

	switch r.Method {
	case "GET":
		accountStateGET(w, r)
	case "PATCH":
		accountStatePATCH(w, r)
	default:
		data.State = "Method should be GET or PATCH"
		json.NewEncoder(w).Encode(data)
	}
}

func accountStateGET(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State     string `json:"state"`
		StateMode bool   `json:"statemode"`
	}

	var read struct {
		Name string `json:"name"`
	}

	read.Name = strings.ToLower(strings.Replace(r.URL.Path, "/api/account/state/", "", 1))

	switch read.Name {
	case "f2f", "fsh", "gch":
		// pass
	default: 
		data.State = "Undefined state"
		json.NewEncoder(w).Encode(data)
		return
	}

	var (
		token  string
	)
	switch {
	case isTokenAuthError(w, r, &token):
		return
	case isLifeTokenError(w, r, token):
		return
	}

	user := settings.Users[token]

	switch read.Name {
	case "f2f":
		data.StateMode = user.State.UsedF2F
	case "fsh":
		data.StateMode = user.State.UsedFSH
	case "gch":
		data.StateMode = user.State.UsedGCH
	}

	json.NewEncoder(w).Encode(data)
}

func accountStatePATCH(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State     string `json:"state"`
		StateMode bool   `json:"statemode"`
	}

	var read struct {
		Name string `json:"name"`
	}

	read.Name = strings.ToLower(strings.Replace(r.URL.Path, "/api/account/state/", "", 1))

	switch read.Name {
	case "f2f", "fsh", "gch":
		// pass
	default: 
		data.State = "Undefined state"
		json.NewEncoder(w).Encode(data)
		return
	}

	var (
		token  string
		client = new(gopeer.Client)
	)
	switch {
	case isTokenAuthError(w, r, &token):
		return
	case isLifeTokenError(w, r, token):
		return
	case isGetClientError(w, r, client, token):
		return
	}

	user := settings.Users[token]
	client = settings.Listener.Clients[user.Hashname]

	switch read.Name {
	case "f2f":
		user.State.UsedF2F = !user.State.UsedF2F
		err := db.SetState(user, &user.State)
		if err != nil {
			user.State.UsedF2F = !client.F2F.Perm
			data.State = "Can't set state 'f2f'"
			json.NewEncoder(w).Encode(data)
			return
		}
		client.F2F.Perm = user.State.UsedF2F
		data.StateMode  = user.State.UsedF2F
	case "fsh":
		user.State.UsedFSH = !user.State.UsedFSH
		err := db.SetState(user, &user.State)
		if err != nil {
			user.State.UsedFSH = !client.Sharing.Perm
			data.State = "Can't set state 'fsh'"
			json.NewEncoder(w).Encode(data)
			return
		}
		client.Sharing.Perm = user.State.UsedFSH
		data.StateMode      = user.State.UsedFSH
	case "gch":
		user.State.UsedGCH = !user.State.UsedGCH
		err := db.SetState(user, &user.State)
		if err != nil {
			user.State.UsedGCH = !user.State.UsedGCH
			data.State = "Can't set state 'gch'"
			json.NewEncoder(w).Encode(data)
			return
		}
		data.StateMode = user.State.UsedGCH
	}

	json.NewEncoder(w).Encode(data)
}
