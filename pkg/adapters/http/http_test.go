package http

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/hidden-lake/pkg/adapters"
	"github.com/number571/hidden-lake/pkg/adapters/http/client"
	"github.com/number571/hidden-lake/pkg/adapters/http/settings"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestHTTPAdapter(t *testing.T) {
	t.Parallel()

	adapter := NewHTTPAdapter(
		NewSettings(&SSettings{
			FAdapterSettings: adapters.NewSettings(&adapters.SSettings{
				FMessageSizeBytes: 8192,
			}),
			FAddress: testutils.TgAddrs[18],
		}),
		cache.NewLRUCache(1024),
		func() []string { return nil },
	)

	adapter.WithHandlers(
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

	go func() { _ = adapter.Run(ctx) }()

	client := client.NewClient(
		client.NewRequester(testutils.TgAddrs[18], &http.Client{Timeout: 5 * time.Second}),
	)

	res, err := client.GetIndex(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	if res != "http-adapter" {
		t.Error("failed get index")
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

	onlines, err := client.GetOnlines(ctx)
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

	netMsg := message.NewMessage(
		message.NewConstructSettings(&message.SConstructSettings{
			FSettings: message.NewSettings(&message.SSettings{}),
		}),
		payload.NewPayload32(0x01, random.NewRandom().GetBytes(8192)),
	)

	if err := client.ProduceMessage(ctx, netMsg); err != nil {
		t.Error(err)
		return
	}
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
