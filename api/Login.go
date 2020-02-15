package api

import (
	"encoding/json"
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/db"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data struct {
		Token    string `json:"token"`
		Hashname string `json:"hashname"`
		State    string `json:"state"`
	}

	if r.Method != "POST" {
		data.State = "Method should be POST"
		json.NewEncoder(w).Encode(data)
		return
	}

	var (
		read = new(userdata)
		user = new(models.User)
	)

	switch {
	case isDecodeError(w, r, read): return
	case isGetUserError(w, r, user, read): return
	}

	token := gopeer.Base64Encode(gopeer.GenerateRandomBytes(20))
	hash := user.Hashname

	if token, ok := settings.Tokens[hash]; ok {
		delete(settings.Users, token)
	}

	settings.Users[token] = user
	settings.Tokens[hash] = token

	client := settings.Listener.NewClient(user.Keys.Private)

	client.SetFriends(user.UsedF2F, db.GetAllFriends(user)...)
	client.SetSharing(true, settings.PATH_ARCHIVE)

	data.Token = token
	data.Hashname = hash
	json.NewEncoder(w).Encode(data)
}
