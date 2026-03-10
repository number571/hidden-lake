package broker

import (
	"testing"
)

func TestMessageBroker(t *testing.T) {
	t.Parallel()

	chanSize := uint64(256)
	subLimit := uint64(256)

	msgBroker := NewDataBroker(chanSize, subLimit)

	_ = msgBroker.Consume("consumer_id")
	msgChan := msgBroker.Consume("consumer_id")

	msg := "hello, world"
	msgBroker.Produce(msg)

	if msgBroker.CountSubscribers() != 1 {
		t.Fatal("count subscribers != 1")
	}

	select {
	case x := <-msgChan:
		if x.(string) != msg {
			t.Fatal("x != msg")
		}
	default:
		t.Fatal("chan is empty")
	}

	// check auto close channel after overflow
	newMsg := "new_message"
	for range chanSize + 1 {
		msgBroker.Produce(newMsg)
	}
	for range msgChan {
		// if channel closed -> cycle has end
	}

	if msgBroker.CountSubscribers() != 0 {
		t.Fatal("count subscribers != 0")
	}

	// update channel
	msgChan = msgBroker.Consume("consumer_id")
	select {
	case <-msgChan:
		t.Fatal("chan is not refreshed")
	default:
	}
	msgBroker.Produce(newMsg)

	select {
	case x := <-msgChan:
		if x.(string) != newMsg {
			t.Fatal("x != newMsg")
		}
	default:
		t.Fatal("chan is empty")
	}
}
