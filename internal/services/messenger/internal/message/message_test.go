package message

import (
	"context"
	"fmt"
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

	if err := msgBroker.Register("consumer_id"); err != nil {
		t.Fatal(err)
	}

	friendName := "Alice"
	msg := dto.NewMessage(true, "hello, world", time.Now())
	msgBroker.Produce(NewMessageContainer(friendName, msg))

	if msgBroker.CountSubscribers() != 1 {
		t.Fatal("count subscribers != 1")
	}

	ctx := context.Background()
	x, err := msgBroker.Consume(ctx, "consumer_id")
	if err != nil {
		t.Fatal(err)
	}

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

	// check auto close channel after overflow
	newMsg := dto.NewMessage(false, "new_message", time.Now())
	for range chanSize + 1 {
		msgBroker.Produce(NewMessageContainer(friendName, newMsg))
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

	msgBroker.Produce(NewMessageContainer(friendName, newMsg))

	x1, err := msgBroker.Consume(ctx, "consumer_id")
	if err != nil {
		t.Fatal(err)
	}
	v1, ok := x1.(IMessageContainer)
	if !ok {
		t.Fatal("invalid type")
	}
	if v1.GetFriend() != friendName {
		t.Fatal("x.GetFriend() != friendName")
	}
	if v1.GetMessage().ToString() != newMsg.ToString() {
		t.Fatal("x.GetMessage().ToString() != newMsg.ToString()")
	}

	for i := range subLimit - 1 {
		if err := msgBroker.Register(fmt.Sprintf("%d", i)); err != nil { // nolint: perfsprint
			t.Fatal(err)
		}
	}
	if err := msgBroker.Register("9999"); err == nil { // nolint: perfsprint
		t.Fatal("success consume with overflow subscribes")
	}
}
