package message

import (
	"testing"
	"time"

	"github.com/number571/hidden-lake/pkg/api/services/messenger/client/dto"
)

func TestMessageBroker(t *testing.T) {
	t.Parallel()

	msgBroker := NewMessageBroker()
	msgChan := msgBroker.Consume("consumer_id")

	friendName := "Alice"
	msg := dto.NewMessage(true, "hello, world", time.Now())
	msgBroker.Produce(friendName, msg)

	select {
	case x := <-msgChan:
		if x.GetFriend() != friendName {
			t.Fatal("x.GetFriend() != friendName")
		}
		if x.GetMessage().ToString() != msg.ToString() {
			t.Fatal("x.GetMessage().ToString() != msg.ToString()")
		}
	default:
		t.Fatal("chan is empty")
	}

	// check auto close channel after overflow
	newMsg := dto.NewMessage(false, "new_message", time.Now())
	for range subscribeChanSize + 1 {
		msgBroker.Produce(friendName, newMsg)
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
	msgBroker.Produce(friendName, newMsg)

	select {
	case x := <-msgChan:
		if x.GetFriend() != friendName {
			t.Fatal("x.GetFriend() != friendName")
		}
		if x.GetMessage().ToString() != newMsg.ToString() {
			t.Fatal("x.GetMessage().ToString() != newMsg.ToString()")
		}
	default:
		t.Fatal("chan is empty")
	}
}
