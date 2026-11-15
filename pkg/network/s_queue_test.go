package network

import (
	"context"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/anonymity/qb/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2/symmetric"
)

func TestSilentQueue(t *testing.T) {
	t.Parallel()

	sQueue := newSilentQueueProcessor(
		queue.NewSettings(&queue.SSettings{
			FMessageConstructSettings: layer1.NewConstructSettings(
				&layer1.SConstructSettings{
					FSettings: layer1.NewSettings(&layer1.SSettings{}),
				},
			),
			FQueuePoolCap: [2]uint64{1, 1}, // must be != 0
			FQueuePeriod:  1,               // must be != 0
			FConsumersCap: 1,               // must be != 0
		}),
		func() layer2.IScheme {
			scheme, _ := symmetric.NewScheme(200)
			return scheme
		}(),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		time.Sleep(100 * time.Millisecond)
	}()

	go func() { _ = sQueue.Run(ctx) }()
	time.Sleep(100 * time.Millisecond)

	if err := sQueue.Run(ctx); err == nil {
		t.Fatal("success double run queue")
	}

	_ = sQueue.GetSettings()
	_ = sQueue.GetScheme()

	msg := []byte("hello, world!")
	if err := sQueue.EnqueueMessage(asymmetric.NewPrivKey().GetPubKey(), msg); err == nil {
		t.Fatal("success enqueue message with another symmetric scheme")
	}

	key := sQueue.GetScheme().GetRandomKey()
	if err := sQueue.EnqueueMessage(key, msg); err != nil {
		t.Fatal(err)
	}

	encMsg := sQueue.DequeueMessage(ctx)
	if encMsg == nil {
		t.Fatal("dequeue message is nil")
	}

	canceledCtx, cancel2 := context.WithCancel(ctx)
	cancel2()
	if m := sQueue.DequeueMessage(canceledCtx); m != nil {
		t.Fatal("success dequeue message with canceled context")
	}

	// message cached into fOutQueue
	_ = sQueue.EnqueueMessage(key, msg)
	_ = sQueue.EnqueueMessage(key, msg)

	if err := sQueue.EnqueueMessage(key, msg); err == nil {
		t.Fatal("success enqueue into overflowed queue")
	}
}
