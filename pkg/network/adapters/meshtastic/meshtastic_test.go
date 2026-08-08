package meshtastic

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/hidden-lake/internal/adapters/meshtastic/pkg/settings"
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

func testRemoveFiles(pPath string) {
	_ = os.Remove(filepath.Join(pPath, settings.CPathTxt))
	_ = os.Remove(filepath.Join(pPath, settings.CPathPy))
	_ = os.RemoveAll(filepath.Join(pPath, settings.CPathVenv))
}

func TestMeshtasticAdapterConfiguration(t *testing.T) {
	t.Parallel()

	path := "./testdata"

	testRemoveFiles(path)
	defer testRemoveFiles(path)

	_ = NewSettings(nil)

	adapterSettings := adapters.NewSettings(&adapters.SSettings{
		FMessageSizeBytes: CLimitMessageSizeBytes,
		FWorkSizeBits:     0,
		FNetworkKey:       "",
	})

	settings := NewSettings(&SSettings{
		FAdapterSettings: adapterSettings,
		FServeSettings: &SServeSettings{
			FPath:         path,
			FAddress:      testutils.TgAddrs[8],
			FDevPath:      "/dev/ttyUSB0",
			FChannel:      1,
			FWatchPeriod:  300 * time.Millisecond,
			FReadTimeout:  2,
			FWriteTimeout: 3,
			FMaxDelayTime: 100 * time.Millisecond,
		},
	})

	meshtasticAdapter := NewMeshtasticAdapter(
		settings,
		cache.NewLRUCache(16),
	)

	if settings.GetPath() != path {
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
	if settings.GetWatchPeriod() != 300*time.Millisecond {
		t.Fatal("got invalid watch_period")
	}
	if settings.GetReadTimeout() != 2 {
		t.Fatal("got invalid read_timeout")
	}
	if settings.GetWriteTimeout() != 3 {
		t.Fatal("got invalid write_timeout")
	}
	if settings.GetMaxDelayTime() != 100*time.Millisecond {
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

	ctx1, cancel1 := context.WithCancel(ctx)
	cancel1()

	if _, err := meshtasticAdapter.Consume(ctx1); err == nil {
		t.Fatal("success consume message with closed context")
	}

	_meshtasticAdapter := meshtasticAdapter.(*sMeshtasticAdapter)
	if _, err := _meshtasticAdapter.getFreePort(); err != nil {
		t.Fatal(err)
	}

	if err := _meshtasticAdapter.createPythonVenv(ctx1); err == nil {
		t.Fatal("success create python venv with closed context")
	}
	if err := _meshtasticAdapter.installPythonRequirements(ctx1); err == nil {
		t.Fatal("success install python requirements with closed context")
	}
	if err := _meshtasticAdapter.runPythonScript(ctx1); err == nil {
		t.Fatal("success run python script with closed context")
	}
	if err := _meshtasticAdapter.closePythonScript(); err == nil {
		t.Fatal("success close python script with undefined service")
	}
	if err := _meshtasticAdapter.closePythonScript(); err == nil {
		t.Fatal("success close not exist python script")
	}

	for range netMessageChanSize {
		if ok := _meshtasticAdapter.pushMessageToChan(msg); !ok {
			t.Fatal("failed push message to chan")
		}
	}

	if ok := _meshtasticAdapter.pushMessageToChan(msg); ok {
		t.Fatal("success push message to overflow chan")
	}
	if _, err := meshtasticAdapter.Consume(ctx); err != nil {
		t.Fatal(err)
	}

	settings1 := NewSettings(&SSettings{
		FAdapterSettings: adapterSettings,
		FServeSettings: &SServeSettings{
			FPath:         path,
			FDevPath:      "/dev/ttyUSB0",
			FChannel:      1,
			FWatchPeriod:  300 * time.Millisecond,
			FReadTimeout:  2,
			FWriteTimeout: 3,
			FMaxDelayTime: 100 * time.Millisecond,
		},
	})

	meshtasticAdapter1 := NewMeshtasticAdapter(
		settings1,
		cache.NewLRUCache(16),
	)

	_meshtasticAdapter1 := meshtasticAdapter1.(*sMeshtasticAdapter)
	if err := _meshtasticAdapter1.runPythonScript(ctx1); err == nil {
		t.Fatal("success run python script with canceled context")
	}
	if err := _meshtasticAdapter.runSubscriber(ctx1); err == nil {
		t.Fatal("success run subscriber with canceled context")
	}

	canceled := make(chan struct{})

	ctx2, cancel2 := context.WithCancel(ctx)
	go func() {
		_ = _meshtasticAdapter.runSubscriber(ctx2)
		canceled <- struct{}{}
	}()

	time.Sleep(time.Second)
	cancel2()
	<-canceled
}
