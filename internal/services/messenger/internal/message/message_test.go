package message

import (
	"testing"
	"time"

	"github.com/number571/hidden-lake/internal/utils/broker"
	"github.com/number571/hidden-lake/pkg/api/services/messenger/client/dto"
)

func TestMessageBroker(t *testing.T) {
	t.Parallel()

	chanSize := uint64(256)
	subLimit := uint64(256)
	msgBroker := broker.NewDataBroker(chanSize, subLimit)

	_ = msgBroker.Consume("consumer_id")
	msgChan := msgBroker.Consume("consumer_id")

	friendName := "Alice"
	msg := dto.NewMessage(true, "hello, world", time.Now())
	msgBroker.Produce(NewMessageContainer(friendName, msg))

	select {
	case x := <-msgChan:
		v, ok := x.(IMessageContainer)
		if !ok {
			t.Fatal("invalid type")
		}
		if v.GetFriend() != friendName {
			t.Fatal("x.GetFriend() != friendName")
		}
		if v.GetMessage().ToString() != msg.ToString() {
			t.Fatal("x.GetMessage().ToString() != msg.ToString()")
		}
	default:
		t.Fatal("chan is empty")
	}

	// check auto close channel after overflow
	newMsg := dto.NewMessage(false, "new_message", time.Now())
	for range chanSize + 1 {
		msgBroker.Produce(NewMessageContainer(friendName, newMsg))
	}
	for range msgChan {
		// if channel closed -> cycle has end
	}

	// update channel
	msgChan = msgBroker.Consume("consumer_id")
	select {
	case <-msgChan:
		t.Fatal("chan is not refreshed")
	default:
	}
	msgBroker.Produce(NewMessageContainer(friendName, newMsg))

	select {
	case x := <-msgChan:
		v, ok := x.(IMessageContainer)
		if !ok {
			t.Fatal("invalid type")
		}
		if v.GetFriend() != friendName {
			t.Fatal("x.GetFriend() != friendName")
		}
		if v.GetMessage().ToString() != newMsg.ToString() {
			t.Fatal("x.GetMessage().ToString() != newMsg.ToString()")
		}
	default:
		t.Fatal("chan is empty")
	}
}
