package handle

import (
	"bytes"
	"fmt"
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/db"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
	"github.com/number571/hiddenlake/utils"
	"golang.org/x/net/websocket"
)

func NewGroupChatMessage(client *gopeer.Client, founder, message string) *models.GroupChat {
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
	return &models.GroupChat{
		Head: models.GroupChatHead{
			Founder: founder,
			Option:  settings.TITLE_GROUPCHAT,
			Sender: models.GroupChatSender{
				Hashname: hashname,
				Public:   gopeer.StringPublic(client.Public()),
			},
		},
		Body: models.GroupChatBody{
			Data: message,
			Desc: models.GroupChatDesc{
				Rand: gopeer.Base64Encode(random),
				Hash: gopeer.Base64Encode(hash),
				Sign: gopeer.Base64Encode(gopeer.Sign(client.Private(), hash)),
			},
		},
	}
}

func getGroupchat(client *gopeer.Client, pack *gopeer.Package) (set string) {
	var (
		glbcht = new(models.GroupChat)
		token  = settings.Tokens[client.Hashname()]
		user   = settings.Users[token]
	)
	gopeer.UnpackJSON([]byte(pack.Body.Data), glbcht)
	if glbcht == nil {
		return
	}
	switch glbcht.Head.Option {
	// add client to chat
	case gopeer.Get("OPTION_GET").(string):
		if !user.State.UsedGCH {
			return
		}
		if _, ok := user.Temp.ChatMap.Owner[pack.From.Sender.Hashname]; ok {
			return
		}
		user.Temp.ChatMap.Owner[pack.From.Sender.Hashname] = true
		redirectAndStore(user, client, "join to chat", pack.From.Sender.Hashname)
		set = string(gopeer.PackJSON([]string{"accept"}))
		return set
	// del client from chat
	case gopeer.Get("OPTION_SET").(string):
		if !user.State.UsedGCH {
			return
		}
		if _, ok := user.Temp.ChatMap.Owner[pack.From.Sender.Hashname]; !ok {
			return
		}
		delete(user.Temp.ChatMap.Owner, pack.From.Sender.Hashname)
		redirectAndStore(user, client, "exit from chat", pack.From.Sender.Hashname)
		set = string(gopeer.PackJSON([]string{"accept"}))
		return set
	// get list of clients chat
	case settings.TITLE_TESTCONN:
		if !user.State.UsedGCH {
			return
		}
		var list []string
		list = append(list, user.Hashname)
		for hash := range user.Temp.ChatMap.Owner {
			list = append(list, hash)
		}
		set = string(gopeer.PackJSON(list))
		return set
	case settings.TITLE_GROUPCHAT:
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
	random := gopeer.Base64Decode(glbcht.Body.Desc.Rand)
	hash := gopeer.HashSum(bytes.Join(
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
	timestamp := utils.CurrentTime()
	if glbcht.Head.Founder == client.Hashname() {
		if !user.State.UsedGCH {
			return
		}
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
					Title:  settings.TITLE_GROUPCHAT,
					Option: gopeer.Get("OPTION_GET").(string),
				},
				Body: gopeer.Body{
					Data: pack.Body.Data,
				},
			})
		}
		db.SetGroupChat(user, &models.Chat{
			Companion: glbcht.Head.Founder,
			Messages: []models.Message{
				models.Message{
					Name: hashname,
					Text: glbcht.Body.Data,
					Time: timestamp,
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
			Time: timestamp,
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
	db.SetGroupChat(user, &models.Chat{
		Companion: glbcht.Head.Founder,
		Messages: []models.Message{
			models.Message{
				Name: hashname,
				Text: glbcht.Body.Data,
				Time: timestamp,
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
		Time: timestamp,
	}
	if user.Session.Socket != nil && user.Session.Option == models.GROUP_OPTION {
		websocket.JSON.Send(user.Session.Socket, wsdata)
	}
	return set
}

func redirectAndStore(user *models.User, client *gopeer.Client, message, sender string) {
	glbcht := NewGroupChatMessage(
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
				Title:  settings.TITLE_GROUPCHAT,
				Option: gopeer.Get("OPTION_GET").(string),
			},
			Body: gopeer.Body{
				Data: string(gopeer.PackJSON(glbcht)),
			},
		})
	}
	timestamp := utils.CurrentTime()
	db.SetGroupChat(user, &models.Chat{
		Companion: client.Hashname(),
		Messages: []models.Message{
			models.Message{
				Name: sender,
				Text: message,
				Time: timestamp,
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
		Time: timestamp,
	}
	if user.Session.Socket != nil && user.Session.Option == models.GROUP_OPTION {
		websocket.JSON.Send(user.Session.Socket, wsdata)
	}
}

func setGroupchat(client *gopeer.Client, pack *gopeer.Package) {
	var (
		list  []string
		token = settings.Tokens[client.Hashname()]
		user  = settings.Users[token]
	)
	gopeer.UnpackJSON([]byte(pack.Body.Data), &list)
	if list == nil {
		return
	}
	if len(list) == 1 && list[0] == "accept" {
		client.Connections[pack.From.Sender.Hashname].Action <- true
		return
	}
	user.Temp.ConnList = []models.Connect{}
	for _, hash := range list {
		user.Temp.ConnList = append(user.Temp.ConnList, models.Connect{Hashname: hash})
	}
	client.Connections[pack.From.Sender.Hashname].Action <- true
}
