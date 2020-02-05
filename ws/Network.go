package ws

import (
	"github.com/number571/hiddenlake/settings"
	"golang.org/x/net/websocket"
)

func Network(ws *websocket.Conn) {
	defer ws.Close()

	var read struct {
		Token string `json:"token"`
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

	websocket.JSON.Receive(ws, &read)
}
