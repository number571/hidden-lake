package handle

import (
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/db"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
	"github.com/number571/hiddenlake/utils"
	"golang.org/x/net/websocket"
)

func getPrivatechat(client *gopeer.Client, pack *gopeer.Package) (set string) {
	var (
		token = settings.Tokens[client.Hashname()]
		hash  = pack.From.Sender.Hashname
		user  = settings.Users[token]
		time  = utils.CurrentTime()
	)

	if !db.InClients(user, hash) {
		db.SetClient(user, &models.Client{
			Hashname:    hash,
			Address:     pack.From.Address,
			Public:      client.Connections[hash].Public(),
			ThrowClient: client.Connections[hash].Throw(),
			Certificate: string(client.Connections[hash].Certificate()),
		})
	}

	message := pack.Body.Data
	if len(message) >= settings.MESSAGE_SIZE {
		return set
	}

	db.SetChat(user, &models.Chat{
		Companion: hash,
		Messages: []models.Message{
			models.Message{
				Name: hash,
				Text: message,
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
			From: hash,
			To:   user.Hashname,
		},
		Text: message,
		Time: time,
	}

	if user.Session.Socket != nil && user.Session.Option == models.PRIVATE_OPTION {
		websocket.JSON.Send(user.Session.Socket, wsdata)
	}

	return set
}

func setPrivatechat(client *gopeer.Client, pack *gopeer.Package) {
	// pass
}
