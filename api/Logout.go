package api

import (
	"encoding/json"
	"github.com/number571/hiddenlake/settings"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data struct {
		State string `json:"state"`
	}

	if r.Method != "POST" {
		data.State = "Method should be POST"
		json.NewEncoder(w).Encode(data)
		return
	}

	var token string
	switch {
	case isTokenAuthError(w, r, &token):
		return
	}

	deleteUserAuth(settings.Users[token])
	json.NewEncoder(w).Encode(data)
}
