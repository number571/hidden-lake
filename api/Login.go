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
	case isDecodeError(w, r, read):
		return
	case isGetUserError(w, r, user, read):
		return
	}

	token := gopeer.Base64Encode(gopeer.GenerateRandomBytes(32))
	hash := user.Hashname

	settings.Mutex.Lock()
	if oldtoken, ok := settings.Tokens[hash]; ok {
		delete(settings.Users, oldtoken)
	}

	settings.Users[token] = user
	settings.Tokens[hash] = token
	settings.Mutex.Unlock()

	client := settings.Listener.NewClient(user.Keys.Private)
	friends := db.GetAllFriends(user)

	client.F2F.Perm = user.State.UsedF2F
	for _, hash := range friends {
		client.F2F.Friends[hash] = true
	}

	client.Sharing.Perm = user.State.UsedFSH
	client.Sharing.Path = settings.PATH_ARCHIVE

	// chat := db.GetGroupChat(user, user.Hashname)
	// if chat != nil && chat.Messages == nil {
	// 	db.SetGroupChat(user, &models.Chat{
	// 		Companion: user.Hashname,
	// 		Messages: []models.Message{
	// 			models.Message{
	// 				Name: user.Hashname,
	// 				Text: "init message",
	// 				Time: utils.CurrentTime(),
	// 			},
	// 		},
	// 	})
	// }

	data.Token = token
	data.Hashname = hash
	json.NewEncoder(w).Encode(data)
}
