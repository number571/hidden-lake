package meshtastic

import (
	"context"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/hidden-lake/pkg/network/adapters"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestMeshtasticAdapter(t *testing.T) {
	t.Parallel()

	adapterSettings := adapters.NewSettings(&adapters.SSettings{
		FMessageSizeBytes: CLimitMessageSizeBytes,
		FWorkSizeBits:     0,
		FNetworkKey:       "",
	})

	meshtasticAdapter := NewMeshtasticAdapter(
		NewSettings(&SSettings{
			FAdapterSettings: adapterSettings,
			FServeSettings: &SServeSettings{
				FPath:         "./testdata",
				FAddress:      testutils.TgAddrs[8],
				FDevPath:      "",
				FChannel:      0,
				FWatchPeriod:  time.Second,
				FReadTimeout:  time.Second,
				FWriteTimeout: time.Second,
				FMaxDelayTime: 1,
			},
		}),
		cache.NewLRUCache(16),
	)

	ctx := context.Background()
	msg := layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: adapterSettings,
		}),
		random.NewRandom().GetBytes(CLimitMessageSizeBytes),
	)
	if err := meshtasticAdapter.Produce(ctx, msg); err == nil {
		t.Fatal("success produce message to undefined service (1)")
	}

	msg2 := layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: adapterSettings,
		}),
		[]byte{1},
	)
	if err := meshtasticAdapter.Produce(ctx, msg2); err == nil {
		t.Fatal("success produce message with invalid size")
	}

	ctx1, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := meshtasticAdapter.Consume(ctx1); err == nil {
		t.Fatal("success consume message with closed context")
	}
}
