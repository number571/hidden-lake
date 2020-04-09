package ws

import (
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
	"golang.org/x/net/websocket"
	"strings"
)

func Network(ws *websocket.Conn) {
	defer ws.Close()

	var read struct {
		Token  string `json:"token"`
		Option string `json:"option"`
	}

	if err := websocket.JSON.Receive(ws, &read); err != nil {
		return
	}

	token := read.Token
	if _, ok := settings.Users[token]; !ok {
		return
	}

	user := settings.Users[token]

	user.Session.Socket = ws
	switch strings.ToLower(read.Option) {
	case "private":
		user.Session.Option = models.PRIVATE_OPTION
	case "group":
		user.Session.Option = models.GROUP_OPTION
	default:
		return
	}

	websocket.JSON.Receive(ws, &read)
}
