package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/cache"
	testutils_gopeer "github.com/number571/go-peer/test/utils"
	"github.com/number571/hidden-lake/pkg/adapters"
	"github.com/number571/hidden-lake/pkg/adapters/http/client"
	"github.com/number571/hidden-lake/pkg/adapters/http/settings"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestHTTPAdapter(t *testing.T) {
	t.Parallel()

	adapterSettings := adapters.NewSettings(&adapters.SSettings{
		FMessageSizeBytes: 8192,
	})

	adapter3 := NewHTTPAdapter(
		NewSettings(&SSettings{
			FAdapterSettings: adapterSettings,
		}),
		cache.NewLRUCache(1),
		func() []string { return nil },
	)

	adapter2 := NewHTTPAdapter(
		NewSettings(&SSettings{
			FAdapterSettings: adapterSettings,
			FAddress:         testutils.TgAddrs[19],
		}),
		cache.NewLRUCache(1024),
		func() []string { return nil },
	)

	adapter1 := NewHTTPAdapter(
		NewSettings(&SSettings{
			FAdapterSettings: adapterSettings,
			FAddress:         testutils.TgAddrs[18],
		}),
		cache.NewLRUCache(1024),
		func() []string { return []string{testutils.TgAddrs[19]} },
	)

	onlines := adapter1.GetOnlines()
	if len(onlines) != 1 || onlines[0] != testutils.TgAddrs[19] {
		t.Error("adapter: get onlines")
		return
	}

	adapter1.WithHandlers(
		NewHandler(
			settings.CHandleIndexPath,
			func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "http-adapter") },
		),
		NewHandler(
			settings.CHandleConfigSettingsPath,
			func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, `{"message_size_bytes":8192}`) },
		),
		NewHandler(
			settings.CHandleConfigConnectsPath,
			func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, `["tcp://abc_1"]`) },
		),
		NewHandler(
			settings.CHandleNetworkOnlinePath,
			func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, `["tcp://abc_2"]`) },
		),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() { _ = adapter3.Run(ctx) }()
	go func() { _ = adapter2.Run(ctx) }()
	go func() { _ = adapter1.Run(ctx) }()

	client := client.NewClient(
		client.NewRequester(testutils.TgAddrs[18], &http.Client{Timeout: 5 * time.Second}),
	)

	err1 := testutils_gopeer.TryN(
		50,
		10*time.Millisecond,
		func() error {
			res, err := client.GetIndex(ctx)
			if err != nil {
				return err
			}
			if res != "http-adapter" {
				t.Error()
				return errors.New("failed get index") // nolint: err113
			}
			return nil
		},
	)
	if err1 != nil {
		t.Error(err1)
		return
	}

	sett, err := client.GetSettings(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	if sett.GetMessageSizeBytes() != 8192 {
		t.Error("invalid settings")
		return
	}

	conns, err := client.GetConnections(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	if len(conns) != 1 || conns[0] != "tcp://abc_1" {
		t.Error("invalid connections")
		return
	}

	onlines, err = client.GetOnlines(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	if len(onlines) != 1 || onlines[0] != "tcp://abc_2" {
		t.Error("invalid onlines")
		return
	}

	if err := client.AddConnection(ctx, "tcp://new"); err != nil {
		t.Error(err)
		return
	}
	if err := client.DelConnection(ctx, "tcp://new"); err != nil {
		t.Error(err)
		return
	}
	if err := client.DelOnline(ctx, "tcp://new"); err != nil {
		t.Error(err)
		return
	}

	msgBytes := []byte("hello, world!")
	msgBytes = append(msgBytes, random.NewRandom().GetBytes(uint64(8192-len(msgBytes)))...)
	netMsg := message.NewMessage(
		message.NewConstructSettings(&message.SConstructSettings{
			FSettings: message.NewSettings(&message.SSettings{}),
		}),
		payload.NewPayload32(0x01, msgBytes),
	)

	if err := client.ProduceMessage(ctx, netMsg); err != nil {
		t.Error(err)
		return
	}
	if err := client.ProduceMessage(ctx, netMsg); err == nil {
		t.Error("success produce message duplicate")
		return
	}

	netMsg2 := message.NewMessage(
		message.NewConstructSettings(&message.SConstructSettings{
			FSettings: message.NewSettings(&message.SSettings{}),
		}),
		payload.NewPayload32(0x01, []byte{1}),
	)
	if err := client.ProduceMessage(ctx, netMsg2); err == nil {
		t.Error("success produce invalid size message")
		return
	}

	if err := testCustomProduceMessage(ctx, http.MethodGet, testutils.TgAddrs[18], netMsg.ToString()); err == nil {
		t.Error("success produce message with invalid method")
		return
	}
	invalidMsg := random.NewRandom().GetString(8192 << 1)
	if err := testCustomProduceMessage(ctx, http.MethodPost, testutils.TgAddrs[18], invalidMsg); err == nil {
		t.Error("success produce message with invalid message")
		return
	}

	msg, err := adapter1.Consume(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.HasPrefix(msg.GetPayload().GetBody(), msgBytes) {
		t.Error("get invalid message bytes")
		return
	}

	if err := adapter1.Produce(ctx, msg); err != nil {
		t.Error(err)
		return
	}
}

func testCustomProduceMessage(ctx context.Context, method, host, msg string) error {
	req, err := http.NewRequestWithContext(ctx, method, "http://"+host, strings.NewReader(msg))
	if err != nil {
		return err
	}
	client := http.Client{Timeout: 5 * time.Second}
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}
	return nil
}

func TestHandler(t *testing.T) {
	t.Parallel()

	path := "/path"
	handler := NewHandler(path, func(_ http.ResponseWriter, _ *http.Request) {})
	if handler.GetPath() != path {
		t.Error("path is invalid")
		return
	}
	_ = handler.GetFunc()
}

func TestSettings(t *testing.T) {
	t.Parallel()

	sett := NewSettings(nil)
	if sett.GetAdapterSettings() == nil {
		t.Error("invalid adapter settings")
		return
	}
}
