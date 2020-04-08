package handle

import (
	"fmt"
	"bytes"
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/db"
	"github.com/number571/hiddenlake/utils"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
	"golang.org/x/net/websocket"
)

func NewGlobalChatMessage(client *gopeer.Client, founder, message string) *models.GlobalChat {
	var (
		hashname = gopeer.HashPublic(client.Public())
		random   = gopeer.GenerateRandomBytes(16)
		hash     = gopeer.HashSum(bytes.Join(
			[][]byte{
				[]byte(hashname),
				[]byte(founder),
				[]byte(message),
				random,
			},
			[]byte{},
		))
	)
	return &models.GlobalChat{
		Head: models.GlobalChatHead{
			Founder: founder,
			Option: settings.TITLE_GLOBALCHAT,
			Sender: models.GlobalChatSender{
				Hashname: hashname,
				Public: gopeer.StringPublic(client.Public()),
			},
		},
		Body: models.GlobalChatBody{
			Data: message,
			Desc: models.GlobalChatDesc{
				Rand: gopeer.Base64Encode(random),
				Hash: gopeer.Base64Encode(hash),
				Sign: gopeer.Base64Encode(gopeer.Sign(client.Private(), hash)),
			},
		},
	}
}

func getGlobalchat(client *gopeer.Client, pack *gopeer.Package) (set string) {
	var (
		glbcht = new(models.GlobalChat)
		token = settings.Tokens[client.Hashname()]
		user  = settings.Users[token]
	)
	gopeer.UnpackJSON([]byte(pack.Body.Data), glbcht)
	if glbcht == nil {
		return
	}
	switch glbcht.Head.Option {
	// add client to chat
	case gopeer.Get("OPTION_GET").(string):
		if _, ok := user.Temp.ChatMap.Owner[pack.From.Sender.Hashname]; ok {
			return
		}
		user.Temp.ChatMap.Owner[pack.From.Sender.Hashname] = true
		redirectAndStore(user, client, "join to chat", pack.From.Sender.Hashname)
		return 
	// del client from chat
	case gopeer.Get("OPTION_SET").(string):
		if _, ok := user.Temp.ChatMap.Owner[pack.From.Sender.Hashname]; !ok {
			return
		}
		delete(user.Temp.ChatMap.Owner, pack.From.Sender.Hashname)
		redirectAndStore(user, client, "exit from chat", pack.From.Sender.Hashname)
		return
	// get list of clients chat
	case settings.TITLE_LOCALCHAT:
		var list []string
		list = append(list, user.Hashname)
		for hash := range user.Temp.ChatMap.Owner {
			list = append(list, hash)
		}
		return string(gopeer.PackJSON(list))
	case settings.TITLE_GLOBALCHAT:
		// pass
	default:
		return
	}
	if len(glbcht.Body.Data) >= settings.MESSAGE_SIZE {
		return
	}
	public := gopeer.ParsePublic(glbcht.Head.Sender.Public)
	if public == nil {
		return
	}
	hashname := gopeer.HashPublic(public)
	if hashname != glbcht.Head.Sender.Hashname {
		return
	}
	random   := gopeer.Base64Decode(glbcht.Body.Desc.Rand)
	hash     := gopeer.HashSum(bytes.Join(
		[][]byte{
			[]byte(hashname),
			[]byte(glbcht.Head.Founder),
			[]byte(glbcht.Body.Data),
			random,
		},
		[]byte{},
	))
	if gopeer.Base64Encode(hash) != glbcht.Body.Desc.Hash {
		return
	}
	if gopeer.Verify(public, hash, gopeer.Base64Decode(glbcht.Body.Desc.Sign)) != nil {
		return
	}
	time := utils.CurrentTime()
	if glbcht.Head.Founder == client.Hashname() {
		if hashname != pack.From.Sender.Hashname {
			return
		}
		if _, ok := user.Temp.ChatMap.Owner[hashname]; !ok {
			return 
		}
		for hash := range user.Temp.ChatMap.Owner {
			if hash == hashname {
				continue
			}
			if !client.InConnections(hash) {
				delete(user.Temp.ChatMap.Owner, hash)
				continue
			}
			dest := client.Destination(hash)
			client.SendTo(dest, &gopeer.Package{
				Head: gopeer.Head{
					Title: settings.TITLE_GLOBALCHAT,
					Option: gopeer.Get("OPTION_GET").(string),
				},
				Body: gopeer.Body{
					Data: pack.Body.Data,
				},
			})
		}
		db.SetGlobalChat(user, &models.Chat{
			Companion: glbcht.Head.Founder,
			Messages: []models.Message{
				models.Message{
					Name: hashname,
					Text: glbcht.Body.Data,
					Time: time,
				},
			},
		})
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
				From: hashname,
				To:   glbcht.Head.Founder,
			},
			Text: glbcht.Body.Data,
			Time: time,
		}
		if user.Session.Socket != nil && user.Session.Option == models.GROUP_OPTION {
			websocket.JSON.Send(user.Session.Socket, wsdata)
		}
		return set
	}
	if _, ok := user.Temp.ChatMap.Member[pack.From.Sender.Hashname]; !ok {
		return
	}
	if glbcht.Head.Founder != pack.From.Sender.Hashname {
		return
	}
	db.SetGlobalChat(user, &models.Chat{
		Companion: glbcht.Head.Founder,
		Messages: []models.Message{
			models.Message{
				Name: hashname,
				Text: glbcht.Body.Data,
				Time: time,
			},
		},
	})
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
			From: hashname,
			To:   glbcht.Head.Founder,
		},
		Text: glbcht.Body.Data,
		Time: time,
	}
	if user.Session.Socket != nil && user.Session.Option == models.GROUP_OPTION {
		websocket.JSON.Send(user.Session.Socket, wsdata)
	}
	return set
}

func redirectAndStore(user *models.User, client *gopeer.Client, message, sender string) {
	glbcht := NewGlobalChatMessage(
		client, 
		client.Hashname(), 
		fmt.Sprintf("%s %s", sender, message),
	)
	for hash := range user.Temp.ChatMap.Owner {
		if !client.InConnections(hash) {
			delete(user.Temp.ChatMap.Owner, hash)
			continue
		}
		dest := client.Destination(hash)
		client.SendTo(dest, &gopeer.Package{
			Head: gopeer.Head{
				Title: settings.TITLE_GLOBALCHAT,
				Option: gopeer.Get("OPTION_GET").(string),
			},
			Body: gopeer.Body{
				Data: string(gopeer.PackJSON(glbcht)),
			},
		})
	}
	db.SetGlobalChat(user, &models.Chat{
		Companion: client.Hashname(),
		Messages: []models.Message{
			models.Message{
				Name: sender,
				Text: message,
				Time: utils.CurrentTime(),
			},
		},
	})
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
			From: sender,
			To:   client.Hashname(),
		},
		Text: message,
		Time: utils.CurrentTime(),
	}
	if user.Session.Socket != nil && user.Session.Option == models.GROUP_OPTION {
		websocket.JSON.Send(user.Session.Socket, wsdata)
	}
}

func setGlobalchat(client *gopeer.Client, pack *gopeer.Package) {
	var (
		list []string
		token = settings.Tokens[client.Hashname()]
		user  = settings.Users[token]
	)
	gopeer.UnpackJSON([]byte(pack.Body.Data), &list)
	if list == nil {
		return
	}
	user.Temp.ConnList = []models.Connect{}
	for _, hash := range list {
		user.Temp.ConnList = append(user.Temp.ConnList, models.Connect{Hashname: hash})
	}
	client.Connections[pack.From.Sender.Hashname].Chans.Action <- true
}
