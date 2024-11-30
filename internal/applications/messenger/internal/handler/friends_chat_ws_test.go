// nolint: goerr113
package handler

import (
	"net/http"
	"testing"
	"time"

	"github.com/number571/hidden-lake/internal/applications/messenger/internal/msgbroker"
	"github.com/number571/hidden-lake/internal/applications/messenger/internal/utils"
	testutils "github.com/number571/hidden-lake/test/utils"
	"golang.org/x/net/websocket"
)

func TestFriendsChatWS(t *testing.T) {
	t.Parallel()

	msgBroker := msgbroker.NewMessageBroker()

	mux := http.NewServeMux()
	mux.Handle("/", websocket.Handler(FriendsChatWS(msgBroker)))

	addr := testutils.TgAddrs[43]
	srv := &http.Server{
		Addr:        addr,
		Handler:     mux,
		ReadTimeout: time.Second,
	}
	defer srv.Close()
	go func() { _ = srv.ListenAndServe() }()

	time.Sleep(200 * time.Millisecond)

	conn, err := websocket.Dial("ws://"+addr, "ws", "http://localhost")
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	subAddr := "abc"
	if err := websocket.JSON.Send(conn, utils.SSubscribe{FAddress: subAddr}); err != nil {
		t.Error(err)
		return
	}

	pMsg := utils.SMessage{
		FFileName:  "file.txt",
		FFileData:  "hello, world!",
		FTimestamp: time.Now().String(),
	}
	msgBroker.Produce(subAddr, pMsg)

	cMsg := utils.SMessage{}
	if err := websocket.JSON.Receive(conn, &cMsg); err != nil {
		t.Error(err)
		return
	}

	if pMsg.FFileName != cMsg.FFileName {
		t.Error(`pMsg.FFileName != cMsg.FFileName`)
		return
	}
	if pMsg.FFileData != cMsg.FFileData {
		t.Error(`pMsg.FFileData != cMsg.FFileData`)
		return
	}
	if pMsg.FTimestamp != cMsg.FTimestamp {
		t.Error(`pMsg.FTimestamp != cMsg.FTimestamp`)
		return
	}
}
