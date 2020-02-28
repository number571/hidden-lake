package handle

import (
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/db"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
	"github.com/number571/hiddenlake/utils"
	"golang.org/x/net/websocket"
)

func Actions(client *gopeer.Client, pack *gopeer.Package) {
	client.HandleAction(settings.TITLE_ARCHIVE, pack, getArchive, setArchive)
	client.HandleAction(settings.TITLE_MESSAGE, pack, getMessage, setMessage)
}

func getArchive(client *gopeer.Client, pack *gopeer.Package) (set string) {
	var (
		token = settings.Tokens[client.Hashname()]
		user  = settings.Users[token]
	)
	if pack.Body.Data == "" {
		return string(gopeer.PackJSON(db.GetAllFiles(user)))
	}
	file := db.GetFile(user, pack.Body.Data)
	if file == nil {
		return ""
	}
	return string(gopeer.PackJSON([]models.File{*file}))
}

func setArchive(client *gopeer.Client, pack *gopeer.Package) {
	var (
		token = settings.Tokens[client.Hashname()]
		user  = settings.Users[token]
	)
	gopeer.UnpackJSON([]byte(pack.Body.Data), &user.Temp.FileList)
	client.Connections[pack.From.Sender.Hashname].Chans.Action <- true
}

func getMessage(client *gopeer.Client, pack *gopeer.Package) (set string) {
	var (
		token    = settings.Tokens[client.Hashname()]
		hash     = pack.From.Sender.Hashname
		user     = settings.Users[token]
		time     = utils.CurrentTime()
	)

	if !db.InClients(user, hash) {
		db.SetClient(user, &models.Client{
			Hashname: hash,
			Address:  pack.From.Address,
			Public:   client.Connections[hash].Public(),
			ThrowClient: client.Connections[hash].Throw(),
			Certificate: string(client.Connections[hash].Certificate()),
		})
	}

	db.SetChat(user, &models.Chat{
		Companion: hash,
		Messages: []models.Message{
			models.Message{
				Name: hash,
				Text: pack.Body.Data,
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
		Text: pack.Body.Data,
		Time: time,
	}

	if user.Session.Socket != nil {
		websocket.JSON.Send(user.Session.Socket, wsdata)
	}

	return set
}

func setMessage(client *gopeer.Client, pack *gopeer.Package) {
	// if package delivered
}
