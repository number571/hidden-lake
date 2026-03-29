package broker

import (
	"context"
	"fmt"
	"testing"
)

func TestMessageBroker(t *testing.T) {
	t.Parallel()

	chanSize := uint64(256)
	subLimit := uint64(32)

	msgBroker := NewDataBroker(chanSize, subLimit)

	if err := msgBroker.Register("consumer_id"); err != nil {
		t.Fatal(err)
	}

	msg := "hello, world"
	msgBroker.Produce(msg)

	if msgBroker.CountSubscribers() != 1 {
		t.Fatal("count subscribers != 1")
	}

	ctx := context.Background()
	x, err := msgBroker.Consume(ctx, "consumer_id")
	if err != nil {
		t.Fatal(err)
	}
	if x.(string) != msg {
		t.Fatal("x != msg")
	}

	// check auto close channel after overflow
	newMsg := "new_message"
	for range chanSize + 1 {
		msgBroker.Produce(newMsg)
	}
	if _, err := msgBroker.Consume(ctx, "consumer_id"); err == nil {
		t.Fatal("success get value from closed consumer")
	}
	if msgBroker.CountSubscribers() != 0 {
		t.Fatal("count subscribers != 0")
	}

	// update channel
	if err := msgBroker.Register("consumer_id"); err != nil {
		t.Fatal(err)
	}

	msgBroker.Produce(newMsg)

	x1, err := msgBroker.Consume(ctx, "consumer_id")
	if err != nil {
		t.Fatal(err)
	}
	if x1.(string) != newMsg {
		t.Fatal("x1 != newMsg")
	}

	for i := range subLimit - 1 {
		if err := msgBroker.Register(fmt.Sprintf("%d", i)); err != nil { // nolint: perfsprint
			t.Fatal(err)
		}
	}
	if err := msgBroker.Register("9999"); err == nil { // nolint: perfsprint
		t.Fatal("success consume with overflow subscribes")
	}

	msgBroker1 := NewDataBroker(2, 1)

	if err := msgBroker1.Register("consumer_id"); err != nil {
		t.Fatal(err)
	}

}
