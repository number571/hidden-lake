package handler

import (
	"github.com/number571/hidden-lake/internal/utils/msgdata"
	"golang.org/x/net/websocket"
)

func ChannelsChatWS(pBroker msgdata.IMessageBroker) func(pWS *websocket.Conn) {
	return func(pWS *websocket.Conn) {
		defer pWS.Close()

		subscribe := new(msgdata.SSubscribe)
		if err := websocket.JSON.Receive(pWS, subscribe); err != nil {
			return
		}

		for {
			msg, ok := pBroker.Consume(subscribe.FAddress)
			if !ok {
				return
			}
			if err := websocket.JSON.Send(pWS, msg); err != nil {
				return
			}
		}
	}
}
