package meshtastic

import (
	"context"
	"fmt"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/hidden-lake/pkg/network/adapters"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestPanicMeshtasticAdapter(t *testing.T) {
	t.Parallel()

	for i := 0; i < 2; i++ {
		testPanicMeshtasticAdapter(t, i)
	}
}

func testPanicMeshtasticAdapter(t *testing.T, n int) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("nothing panics")
		}
	}()
	switch n {
	case 0:
		_ = NewMeshtasticAdapter(
			NewSettings(&SSettings{
				FAdapterSettings: adapters.NewSettings(&adapters.SSettings{
					FMessageSizeBytes: 4096,
				}),
			}),
			cache.NewLRUCache(16),
		)
	case 1:
		_ = NewMeshtasticAdapter(
			NewSettings(&SSettings{
				FAdapterSettings: adapters.NewSettings(&adapters.SSettings{
					FMessageSizeBytes: CLimitMessageSizeBytes,
					FWorkSizeBits:     1,
				}),
			}),
			cache.NewLRUCache(16),
		)
	}
}

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestMeshtasticAdapter(t *testing.T) {
	t.Parallel()

	_ = NewSettings(nil)

	adapterSettings := adapters.NewSettings(&adapters.SSettings{
		FMessageSizeBytes: CLimitMessageSizeBytes,
		FWorkSizeBits:     0,
		FNetworkKey:       "",
	})

	settings := NewSettings(&SSettings{
		FAdapterSettings: adapterSettings,
		FServeSettings: &SServeSettings{
			FPath:         "./testdata",
			FAddress:      testutils.TgAddrs[8],
			FDevPath:      "/dev/ttyUSB0",
			FChannel:      1,
			FWatchPeriod:  1,
			FReadTimeout:  2,
			FWriteTimeout: 3,
			FMaxDelayTime: 4,
		},
	})

	meshtasticAdapter := NewMeshtasticAdapter(
		settings,
		cache.NewLRUCache(16),
	)

	if settings.GetPath() != "./testdata" {
		t.Fatal("got invalid path")
	}
	if settings.GetAddress() != testutils.TgAddrs[8] {
		t.Fatal("got invalid address")
	}
	if settings.GetDevPath() != "/dev/ttyUSB0" {
		t.Fatal("got invalid devpath")
	}
	if settings.GetChannel() != 1 {
		t.Fatal("got invalid channel")
	}
	if settings.GetWatchPeriod() != 1 {
		t.Fatal("got invalid watch_period")
	}
	if settings.GetReadTimeout() != 2 {
		t.Fatal("got invalid read_timeout")
	}
	if settings.GetWriteTimeout() != 3 {
		t.Fatal("got invalid write_timeout")
	}
	if settings.GetMaxDelayTime() != 4 {
		fmt.Println(settings.GetMaxDelayTime())
		t.Fatal("got invalid max_delay_time")
	}

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

	_meshtasticAdapter := meshtasticAdapter.(*sMeshtasticAdapter)
	if _, err := _meshtasticAdapter.getFreePort(); err != nil {
		t.Fatal(err)
	}

	if ok := _meshtasticAdapter.pushMessageToChan(msg); !ok {
		t.Fatal("failed push message to chan")
	}
	if _, err := meshtasticAdapter.Consume(ctx); err != nil {
		t.Fatal(err)
	}

	if err := _meshtasticAdapter.closePythonScript(); err == nil {
		t.Fatal("success close not exist python script")
	}
}
