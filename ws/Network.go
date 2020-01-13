package ws

import (
	"../settings"
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

	for {
		websocket.JSON.Receive(ws, &read)
		return
	}
}
