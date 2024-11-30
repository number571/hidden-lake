package msgbroker

import (
	"html/template"
	"testing"
	"time"

	"github.com/number571/hidden-lake/internal/applications/messenger/internal/utils"
)

func TestMessageBroker(t *testing.T) {
	t.Parallel()

	addr := "address"
	msgData := template.HTML("msg_data")

	msgReceiver := NewMessageBroker()

	go func() {
		time.Sleep(100 * time.Millisecond)
		msgReceiver.Produce(addr, utils.SMessage{FTextData: msgData})
	}()

	msg, ok := msgReceiver.Consume(addr)
	if !ok {
		t.Error("got not ok recv")
		return
	}

	if msg.FTextData != msgData {
		t.Error("msg.FTextData != msgData")
		return
	}
}
