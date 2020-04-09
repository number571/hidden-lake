package api

import (
	"encoding/json"
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/db"
	"github.com/number571/hiddenlake/handle"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
	"github.com/number571/hiddenlake/utils"
	"golang.org/x/net/websocket"
	"net/http"
	"strings"
	"time"
)

type chatlist struct {
	Hashname  string `json:"hashname"`
	Connected bool   `json:"connected"`
}

func NetworkChatGlobal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data struct {
		State string `json:"state"`
	}

	switch r.Method {
	case "GET":
		networkChatGlobalGET(w, r)
	case "POST":
		networkChatGlobalPOST(w, r)
	case "PATCH":
		networkChatGlobalPATCH(w, r)
	case "DELETE":
		networkChatGlobalDELETE(w, r)
	default:
		data.State = "Method should be GET, POST or DELETE"
		json.NewEncoder(w).Encode(data)
	}
}

func networkChatGlobalGET(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State string       `json:"state"`
		List  []chatlist   `json:"list"`
		Chat  *models.Chat `json:"chat"`
	}

	var read struct {
		Hashname string `json:"hashname"`
	}

	read.Hashname = strings.Replace(r.URL.Path, "/api/network/chat/global/", "", 1)

	var (
		client = new(gopeer.Client)
		token  string
		ok     bool
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

	switch read.Hashname {
	case "", "null", "undefined":
		list := db.GetGlobalChatFounders(user)
		for _, hash := range list {
			if hash == user.Hashname {
				ok = true
			} else {
				_, ok = user.Temp.ChatMap.Member[hash]
			}
			data.List = append(data.List, chatlist{
				Hashname:  hash,
				Connected: ok,
			})
		}
	default:
		data.Chat = db.GetGlobalChat(user, read.Hashname)
		if read.Hashname == user.Hashname {
			data.List = append(data.List, chatlist{
				Hashname: user.Hashname,
			})
			for hash := range user.Temp.ChatMap.Owner {
				data.List = append(data.List, chatlist{
					Hashname: hash,
				})
			}
		} else if client.InConnections(read.Hashname) {
			dest := client.Destination(read.Hashname)
			client.SendTo(dest, &gopeer.Package{
				Head: gopeer.Head{
					Title:  settings.TITLE_GLOBALCHAT,
					Option: gopeer.Get("OPTION_GET").(string),
				},
				Body: gopeer.Body{
					Data: string(gopeer.PackJSON(models.GlobalChat{
						Head: models.GlobalChatHead{
							Founder: read.Hashname,
							Option:  settings.TITLE_LOCALCHAT,
						},
					})),
				},
			})
			select {
			case <-client.Connections[read.Hashname].Chans.Action:
				for _, conn := range user.Temp.ConnList {
					data.List = append(data.List, chatlist{
						Hashname: conn.Hashname,
					})
				}
			case <-time.After(time.Duration(gopeer.Get("WAITING_TIME").(uint8)) * time.Second):
				// pass
			}
		}
	}

	json.NewEncoder(w).Encode(data)
}

func networkChatGlobalPOST(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State string `json:"state"`
	}

	var read struct {
		Hashname string `json:"hashname"`
		Message  string `json:"message"`
	}

	var (
		client = new(gopeer.Client)
		token  string
	)
	switch {
	case isTokenAuthError(w, r, &token):
		return
	case isLifeTokenError(w, r, token):
		return
	case isDecodeError(w, r, &read):
		return
	case isGetClientError(w, r, client, token):
		return
	}

	user := settings.Users[token]
	message := strings.Replace(read.Message, "\n", " ", -1)
	if len(message) >= settings.MESSAGE_SIZE {
		data.State = "Message length >= maximum size"
		json.NewEncoder(w).Encode(data)
		return
	}

	glbcht := handle.NewGlobalChatMessage(client, read.Hashname, message)
	if read.Hashname == user.Hashname {
		for hash := range user.Temp.ChatMap.Owner {
			if !client.InConnections(hash) {
				delete(user.Temp.ChatMap.Owner, hash)
				continue
			}
			dest := client.Destination(hash)
			client.SendTo(dest, &gopeer.Package{
				Head: gopeer.Head{
					Title:  settings.TITLE_GLOBALCHAT,
					Option: gopeer.Get("OPTION_GET").(string),
				},
				Body: gopeer.Body{
					Data: string(gopeer.PackJSON(glbcht)),
				},
			})
		}
	} else {
		switch {
		case isNotInConnectionsError(w, r, client, read.Hashname):
			return
		}
		dest := client.Destination(read.Hashname)
		_, err := client.SendTo(dest, &gopeer.Package{
			Head: gopeer.Head{
				Title:  settings.TITLE_GLOBALCHAT,
				Option: gopeer.Get("OPTION_GET").(string),
			},
			Body: gopeer.Body{
				Data: string(gopeer.PackJSON(glbcht)),
			},
		})
		if err != nil {
			data.State = "User can't receive message"
			json.NewEncoder(w).Encode(data)
			return
		}
	}

	time := utils.CurrentTime()
	err := db.SetGlobalChat(user, &models.Chat{
		Companion: read.Hashname,
		Messages: []models.Message{
			models.Message{
				Name: client.Hashname(),
				Text: message,
				Time: time,
			},
		},
	})
	if err != nil {
		data.State = "Set chat error"
		json.NewEncoder(w).Encode(data)
		return
	}

	var wsdata = struct {
		Comp struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"comp"`
		Text string `json:"text"`
		Time string `json:"time"`
	}{
		Comp: struct {
			From string `json:"from"`
			To   string `json:"to"`
		}{
			From: client.Hashname(),
			To:   read.Hashname,
		},
		Text: message,
		Time: time,
	}
	if user.Session.Socket != nil && user.Session.Option == models.GROUP_OPTION {
		websocket.JSON.Send(user.Session.Socket, wsdata)
	}

	json.NewEncoder(w).Encode(data)
}

func networkChatGlobalPATCH(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State string `json:"state"`
	}

	var read struct {
		Hashname string `json:"hashname"`
		Option   string `json:"option"`
	}

	var (
		client = new(gopeer.Client)
		token  string
		err    error
	)
	switch {
	case isTokenAuthError(w, r, &token):
		return
	case isLifeTokenError(w, r, token):
		return
	case isDecodeError(w, r, &read):
		return
	case isGetClientError(w, r, client, token):
		return
	}

	user := settings.Users[token]
	if user.Hashname == read.Hashname {
		data.State = "Can't join/exit to/from own chat"
		json.NewEncoder(w).Encode(data)
		return
	}

	time := utils.CurrentTime()
	switch strings.ToLower(read.Option) {
	case "join":
		if !client.InConnections(read.Hashname) {
			clientData := db.GetClient(user, read.Hashname)
			if clientData == nil {
				data.State = "Client undefined"
				json.NewEncoder(w).Encode(data)
				return
			}
			dest := &gopeer.Destination{
				Address:     clientData.Address,
				Certificate: []byte(clientData.Certificate),
				Public:      clientData.ThrowClient,
				Receiver:    clientData.Public,
			}
			err = client.Connect(dest)
			if err != nil {
				data.State = "Connect error"
				json.NewEncoder(w).Encode(data)
				return
			}
		}
		dest := client.Destination(read.Hashname)
		_, err = client.SendTo(dest, &gopeer.Package{
			Head: gopeer.Head{
				Title:  settings.TITLE_GLOBALCHAT,
				Option: gopeer.Get("OPTION_GET").(string),
			},
			Body: gopeer.Body{
				Data: string(gopeer.PackJSON(models.GlobalChat{
					Head: models.GlobalChatHead{
						Founder: read.Hashname,
						Option:  gopeer.Get("OPTION_GET").(string),
					},
				})),
			},
		})
		if err != nil {
			data.State = "User can't receive message"
			json.NewEncoder(w).Encode(data)
			return
		}
		user.Temp.ChatMap.Member[read.Hashname] = true
		err = db.SetGlobalChat(user, &models.Chat{
			Companion: read.Hashname,
			Messages: []models.Message{
				models.Message{
					Name: client.Hashname(),
					Text: "join to chat",
					Time: time,
				},
			},
		})
		if err != nil {
			data.State = "Set chat error"
			json.NewEncoder(w).Encode(data)
			return
		}
	case "exit":
		delete(user.Temp.ChatMap.Member, read.Hashname)
		switch {
		case isNotInConnectionsError(w, r, client, read.Hashname):
			return
		}
		dest := client.Destination(read.Hashname)
		_, err = client.SendTo(dest, &gopeer.Package{
			Head: gopeer.Head{
				Title:  settings.TITLE_GLOBALCHAT,
				Option: gopeer.Get("OPTION_GET").(string),
			},
			Body: gopeer.Body{
				Data: string(gopeer.PackJSON(models.GlobalChat{
					Head: models.GlobalChatHead{
						Founder: read.Hashname,
						Option:  gopeer.Get("OPTION_SET").(string),
					},
				})),
			},
		})
		if err != nil {
			data.State = "User can't receive message"
			json.NewEncoder(w).Encode(data)
			return
		}
		err = db.SetGlobalChat(user, &models.Chat{
			Companion: read.Hashname,
			Messages: []models.Message{
				models.Message{
					Name: client.Hashname(),
					Text: "exit from chat",
					Time: time,
				},
			},
		})
		if err != nil {
			data.State = "Set chat error"
			json.NewEncoder(w).Encode(data)
			return
		}
	default:
		data.State = "Option undefined"
		json.NewEncoder(w).Encode(data)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func networkChatGlobalDELETE(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State string `json:"state"`
	}

	var (
		read  = new(userdata)
		user  = new(models.User)
		token string
	)
	switch {
	case isTokenAuthError(w, r, &token):
		return
	case isLifeTokenError(w, r, token):
		return
	case isDecodeError(w, r, read):
		return
	case isGetUserError(w, r, user, read):
		return
	case isCheckUserError(w, r, user, token):
		return
	}

	switch strings.ToLower(read.PasswordRepeat) {
	case "clear":
		err := db.ClearGlobalChat(user, read.Hashname)
		if err != nil {
			data.State = "Clear chat error"
			json.NewEncoder(w).Encode(data)
			return
		}
	case "delete":
		err := db.DeleteGlobalChat(settings.Users[token], read.Hashname)
		if err != nil {
			data.State = "Delete chat error"
			json.NewEncoder(w).Encode(data)
			return
		}
	default:
		data.State = "Option undefined"
		json.NewEncoder(w).Encode(data)
		return
	}

	json.NewEncoder(w).Encode(data)
}
