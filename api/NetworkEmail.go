package api

import (
	"strings"
	"net/http"
	"encoding/json"
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/db"
	"github.com/number571/hiddenlake/utils"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/handle"
	"github.com/number571/hiddenlake/settings"
)

func NetworkEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data struct {
		State string `json:"state"`
	}

	switch r.Method {
	case "GET":
		networkEmailGET(w, r)
	case "POST":
		networkEmailPOST(w, r)
	case "PATCH":
		networkEmailPATCH(w, r)
	case "DELETE":
		networkEmailDELETE(w, r)
	default:
		data.State = "Method should be GET, POST, PATCH or DELETE"
		json.NewEncoder(w).Encode(data)
	}
}

// Read email / emails.
func networkEmailGET(w http.ResponseWriter, r *http.Request) {
	hash := strings.Replace(r.URL.Path, "/api/network/email/", "", 1)
	switch hash {
	case "", "null", "undefined":
		allEmailsGET(w, r)
		return
	default:
		oneEmailGET(w, r, hash)
		return
	}
}

func oneEmailGET(w http.ResponseWriter, r *http.Request, hash string) {
	var data struct {
		State   string       `json:"state"`
		Email   models.Email `json:"email"`
	}

	var token string
	switch {
	case isTokenAuthError(w, r, &token): return
	case isLifeTokenError(w, r, token): return
	}

	user := settings.Users[token]
	email := db.GetEmail(user, hash)
	if email == nil {
		data.State = "Error read email"
		json.NewEncoder(w).Encode(data)
		return
	}

	data.Email = *email
	json.NewEncoder(w).Encode(data)
}

func allEmailsGET(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State   string         `json:"state"`
		Emails  []models.Email `json:"emails"`
	}

	var token string
	switch {
	case isTokenAuthError(w, r, &token): return
	case isLifeTokenError(w, r, token): return
	}

	user := settings.Users[token]
	data.Emails = db.GetAllEmails(user)
	
	json.NewEncoder(w).Encode(data)
}

// Update emails.
func networkEmailPATCH(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State string `json:"state"`
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

	packReq := &gopeer.Package{
		Head: gopeer.Head{
			Title: settings.TITLE_EMAIL,
			Option: gopeer.Get("OPTION_GET").(string),
		},
		Body: gopeer.Body{
			Data: string(gopeer.PackJSON(&models.EmailType{
				Head: models.EmailHead{
					Receiver: client.Hashname(),
					Session: gopeer.Get("OPTION_GET").(string),
				},
			})),
		},
	}

	for hash := range client.Connections {
		dest := client.Destination(hash)
		client.SendTo(dest, packReq)
	}

	json.NewEncoder(w).Encode(data)
}

// Send email to another client.
func networkEmailPOST(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State string `json:"state"`
	}

	var read struct {
		PublicKey string `json:"public_key"`
		Message   string `json:"message"`
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

	public := gopeer.ParsePublic(read.PublicKey)
	if public == nil {
		data.State = "Error decode public key"
		json.NewEncoder(w).Encode(data)
		return
	}

	var (
		user  = settings.Users[token]
		hash  = gopeer.HashPublic(public)
		email = handle.NewEmail(client, public, read.Message)
	)

	if client.Hashname() == hash {
		email.Body.Data = read.Message
		err := db.SetEmail(user, models.IsPermEmail, &models.Email{
			Info: models.EmailInfo{
				Incoming: false,
				Time: utils.CurrentTime(),
			},
			Email: *email,
		})
		if err != nil {
			data.State = "Set email error"
			json.NewEncoder(w).Encode(data)
			return
		}
		json.NewEncoder(w).Encode(data)
		return
	}

	if client.InConnections(hash) {
		dest := client.Destination(hash)
		_, err := client.SendTo(dest, &gopeer.Package{
			Head: gopeer.Head{
				Title: settings.TITLE_EMAIL,
				Option: gopeer.Get("OPTION_GET").(string),
			},
			Body: gopeer.Body{
				Data: string(gopeer.PackJSON(email)),
			},
		})
		if err != nil {
			data.State = "User can't receive email"
			json.NewEncoder(w).Encode(data)
			return
		} 
		email.Body.Data = read.Message
		err = db.SetEmail(user, models.IsPermEmail, &models.Email{
			Info: models.EmailInfo{
				Incoming: false,
				Time: utils.CurrentTime(),
			},
			Email: *email,
		})
		if err != nil {
			data.State = "Set email error"
			json.NewEncoder(w).Encode(data)
			return
		}
		json.NewEncoder(w).Encode(data)
		return
	}

	packEmail := string(gopeer.PackJSON(email))
	for hash := range client.Connections {
		dest := client.Destination(hash)
		client.SendTo(dest, &gopeer.Package{
			Head: gopeer.Head{
				Title: settings.TITLE_EMAIL,
				Option: gopeer.Get("OPTION_GET").(string),
			},
			Body: gopeer.Body{
				Data: packEmail,
			},
		})
	}

	email.Body.Data = read.Message
	err := db.SetEmail(user, models.IsPermEmail, &models.Email{
		Info: models.EmailInfo{
			Incoming: false,
			Time: utils.CurrentTime(),
		},
		Email: *email,
	})
	if err != nil {
		data.State = "Set email error"
		json.NewEncoder(w).Encode(data)
		return
	}

	json.NewEncoder(w).Encode(data)
}

// Delete email message.
func networkEmailDELETE(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State string `json:"state"`
	}
	
	var read struct {
		Emailhash string `json:"emailhash"`
	}

	var token string
	switch {
	case isTokenAuthError(w, r, &token): return
	case isLifeTokenError(w, r, token): return
	case isDecodeError(w, r, &read): return
	}

	user := settings.Users[token]
	err := db.DeleteEmail(user, read.Emailhash)
	if err != nil {
		data.State = "Email already deleted"
		json.NewEncoder(w).Encode(data)
		return
	}

	json.NewEncoder(w).Encode(data)
}
