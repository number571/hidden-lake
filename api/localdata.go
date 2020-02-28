package api

import (
	"strings"
	"net/http"
	"encoding/json"
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/db"
	"github.com/number571/hiddenlake/utils"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

type userdata struct {
	Hashname       string `json:"hashname"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"password_repeat"`
	PrivateKey     string `json:"private_key"`
}

func deleteUserAuth(user *models.User) {
	hash := user.Hashname
	token := settings.Tokens[hash]
	delete(settings.Listener.Clients, hash)
	delete(settings.Tokens, hash)
	delete(settings.Users, token)
}

func isDecodeError(w http.ResponseWriter, r *http.Request, read interface{}) bool {
	var data struct {
		State string `json:"state"`
	}
	err := json.NewDecoder(r.Body).Decode(read)
	if err != nil {
		data.State = "Error decode json format"
		json.NewEncoder(w).Encode(data)
		return true 
	}
	return false
}

func isTokenAuthError(w http.ResponseWriter, r *http.Request, token *string) bool {
	var data struct {
		State string `json:"state"`
	}
	*token = r.Header.Get("Authorization")
	*token = strings.Replace(*token, "Bearer ", "", 1)
	if _, ok := settings.Users[*token]; !ok {
		data.State = "Tokened user undefined"
		json.NewEncoder(w).Encode(data)
		return true 
	}
	return false
}

func isLifeTokenError(w http.ResponseWriter, r *http.Request, token string) bool {
	var data struct {
		State string `json:"state"`
	}
	err := settings.CheckLifetimeToken(token)
	if err != nil {
		data.State = "Token lifetime is over"
		json.NewEncoder(w).Encode(data)
		return true
	}
	settings.Users[token].Session.Time = utils.CurrentTime()
	return false
}

func isGetUserError(w http.ResponseWriter, r *http.Request, user *models.User, read *userdata) bool {
	var data struct {
		State string `json:"state"`
	}
	us := db.GetUser(read.Username, read.Password)
	if us == nil {
		data.State = "User undefined"
		json.NewEncoder(w).Encode(data)
		return true
	}
	us = db.GetState(us)
	if us == nil {
		data.State = "Get user state error"
		json.NewEncoder(w).Encode(data)
		return true
	}
	*user = *us
	return false
}

func isGetClientError(w http.ResponseWriter, r *http.Request, client *gopeer.Client, token string) bool {
	var data struct {
		State string `json:"state"`
	}
	hash := settings.Users[token].Hashname
	cl, ok := settings.Listener.Clients[hash]
	if !ok {
		data.State = "Current client is not exist"
		json.NewEncoder(w).Encode(data)
		return true
	}
	*client = *cl
	return false
}

func isNotInConnectionsError(w http.ResponseWriter, r *http.Request, client *gopeer.Client, hashname string) bool {
	var data struct {
		State string `json:"state"`
	}
	if !client.InConnections(hashname) {
		data.State = "User is not connected"
		json.NewEncoder(w).Encode(data)
		return true
	}
	return false
}
