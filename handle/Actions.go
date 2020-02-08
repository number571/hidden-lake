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
	client.HandleAction(settings.TITLE_CONNLIST, pack, getConnlist, setConnlist)
	client.HandleAction(settings.TITLE_ARCHIVE, pack, getArchive, setArchive)
	client.HandleAction(settings.TITLE_MESSAGE, pack, getMessage, setMessage)
}

func getConnlist(client *gopeer.Client, pack *gopeer.Package) (set string) {
	var (
		templist []models.Connect
	)
	if pack.From.Sender.Hashname != pack.From.Hashname {
		return ""
	}
	for hash, cl := range client.Connections {
		if hash == client.Hashname || hash == pack.From.Sender.Hashname {
			continue
		}
		pub1 := gopeer.StringPublic(cl.Public)
		pub2 := gopeer.StringPublic(cl.PublicRecv)
		if pub1 != pub2 {
			continue
		}
		templist = append(templist, models.Connect{
			Hashname: hash,
			PublicKey: pub1,
		})
	}
	return string(gopeer.PackJSON(templist))
}

func setConnlist(client *gopeer.Client, pack *gopeer.Package) {
	var (
		templist []models.Connect
		token = settings.Tokens[client.Hashname]
	)
	gopeer.UnpackJSON([]byte(pack.Body.Data), &templist)
	for i, cl := range templist {
		if cl.Hashname == client.Hashname {
			continue
		}
		if client.InConnections(cl.Hashname) {
			templist[i].Connected = true
		}
	}
	settings.Users[token].Temp.ConnList = templist
}

func getArchive(client *gopeer.Client, pack *gopeer.Package) (set string) {
	var (
		token = settings.Tokens[client.Hashname]
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
		token = settings.Tokens[client.Hashname]
		user  = settings.Users[token]
	)
	gopeer.UnpackJSON([]byte(pack.Body.Data), &user.Temp.FileList)
}

func getMessage(client *gopeer.Client, pack *gopeer.Package) (set string) {
	var (
		hashname = pack.From.Sender.Hashname
		token    = settings.Tokens[client.Hashname]
		user     = settings.Users[token]
		time     = utils.CurrentTime()
	)

	if !db.InClients(user, hashname) {
		db.SetClient(user, &models.Client{
			Hashname: hashname,
			Address:  pack.From.Address,
			Public:   client.Connections[hashname].Public,
		})
	}

	db.SetChat(user, &models.Chat{
		Companion: hashname,
		Messages: []models.Message{
			models.Message{
				Name: hashname,
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
			From: hashname,
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
