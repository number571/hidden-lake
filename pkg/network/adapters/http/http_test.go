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
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/cache"
	testutils_gopeer "github.com/number571/go-peer/test/utils"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/pkg/network/adapters"
	testutils "github.com/number571/hidden-lake/test/utils"

	"github.com/number571/hidden-lake/internal/adapters/http/pkg/settings"
	"github.com/number571/hidden-lake/pkg/api/adapters/http/client"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestHTTPAdapter(t *testing.T) { // nolint: gocyclo, maintidx
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
			FServeSettings: &SServeSettings{
				FAddress: testutils.TgAddrs[19],
			},
		}),
		cache.NewLRUCache(1024),
		func() []string { return nil },
	)

	adapter1 := NewHTTPAdapter(
		NewSettings(&SSettings{
			FAdapterSettings: adapterSettings,
			FServeSettings: &SServeSettings{
				FAddress: testutils.TgAddrs[18],
			},
		}),
		cache.NewLRUCache(1024),
		func() []string { return []string{testutils.TgAddrs[19]} },
	)

	onlines := adapter1.GetOnlines()
	if len(onlines) != 1 || onlines[0] != testutils.TgAddrs[19] {
		t.Fatal("adapter: get onlines")
	}

	adapter1.WithHandlers(map[string]http.HandlerFunc{
		settings.CHandleIndexPath: func(w http.ResponseWriter, _ *http.Request) {
			_, _ = fmt.Fprint(w, "http-adapter")
		},
		settings.CHandleConfigSettingsPath: func(w http.ResponseWriter, _ *http.Request) {
			_, _ = fmt.Fprint(w, `{"message_size_bytes":8192}`)
		},
		settings.CHandleConfigConnectsPath: func(w http.ResponseWriter, _ *http.Request) {
			_, _ = fmt.Fprint(w, `["tcp://abc_1"]`)
		},
		settings.CHandleNetworkOnlinePath: func(w http.ResponseWriter, _ *http.Request) {
			_, _ = fmt.Fprint(w, `["tcp://abc_2"]`)
		},
	})

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
				return errors.New("failed get index") // nolint: err113
			}
			return nil
		},
	)
	if err1 != nil {
		t.Fatal(err1)
	}

	conns, err := client.GetConnections(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(conns) != 1 || conns[0] != "tcp://abc_1" {
		t.Fatal("invalid connections")
	}

	onlines, err = client.GetOnlines(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(onlines) != 1 || onlines[0] != "tcp://abc_2" {
		t.Fatal("invalid onlines")
	}

	if err := client.AddConnection(ctx, "tcp://new"); err != nil {
		t.Fatal(err)
	}
	if err := client.DelConnection(ctx, "tcp://new"); err != nil {
		t.Fatal(err)
	}
	if err := client.DelOnline(ctx, "tcp://new"); err != nil {
		t.Fatal(err)
	}

	msgBytes := []byte("hello, world!")
	msgBytes = append(msgBytes, random.NewRandom().GetBytes(uint64(8192-len(msgBytes)))...) //nolint:gosec
	netMsg := layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{}),
		}),
		payload.NewPayload32(build.GetSettings().FProtoMask.FNetwork, msgBytes),
	)

	if err := client.ProduceMessage(ctx, netMsg); err != nil {
		t.Fatal(err)
	}

	sCtx1, cancel1 := context.WithTimeout(ctx, time.Second)
	defer cancel1()

	req1, err := http.NewRequestWithContext(
		sCtx1,
		http.MethodGet,
		"http://"+testutils.TgAddrs[18]+settings.CHandleNetworkAdapterPath,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	rsp1, err := http.DefaultClient.Do(req1)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rsp1.Body.Close() }()

	if rsp1.StatusCode != http.StatusMethodNotAllowed {
		t.Fatal("got another status code (1)")
	}

	sCtx2, cancel2 := context.WithTimeout(ctx, time.Second)
	defer cancel2()

	req2, err := http.NewRequestWithContext(
		sCtx2,
		http.MethodPost,
		"http://"+testutils.TgAddrs[18]+settings.CHandleNetworkAdapterPath,
		bytes.NewBufferString(encoding.HexEncode([]byte{1})),
	)
	if err != nil {
		t.Fatal(err)
	}

	rsp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rsp2.Body.Close() }()

	if rsp2.StatusCode != http.StatusBadRequest {
		t.Fatal("got another status code (2)")
	}

	sCtx3, cancel3 := context.WithTimeout(ctx, time.Second)
	defer cancel3()

	size := adapterSettings.GetMessageSizeBytes() + layer1.CMessageHeadSize
	req3, err := http.NewRequestWithContext(
		sCtx3,
		http.MethodPost,
		"http://"+testutils.TgAddrs[18]+settings.CHandleNetworkAdapterPath,
		bytes.NewBufferString(encoding.HexEncode(random.NewRandom().GetBytes(size))),
	)
	if err != nil {
		t.Fatal(err)
	}

	rsp3, err := http.DefaultClient.Do(req3)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rsp3.Body.Close() }()

	if rsp3.StatusCode != http.StatusBadRequest {
		t.Fatal("got another status code (3)")
	}

	// duplicate produce is ok (cache checker)
	if err := client.ProduceMessage(ctx, netMsg); err != nil {
		t.Fatal(err)
	}

	netMsg2 := layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{}),
		}),
		payload.NewPayload32(0x01, []byte{1}),
	)
	if err := client.ProduceMessage(ctx, netMsg2); err == nil {
		t.Fatal("success produce invalid size message")
	}

	msgBytesX := []byte("hello, world!")
	msgBytesX = append(msgBytesX, random.NewRandom().GetBytes(uint64(8192-len(msgBytesX)))...) //nolint:gosec
	netMsgX := layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{}),
		}),
		payload.NewPayload32(999, msgBytesX),
	)
	if err := client.ProduceMessage(ctx, netMsgX); err == nil {
		t.Fatal("success produce invalid proto mask")
	}

	if err := testCustomProduceMessage(ctx, http.MethodGet, testutils.TgAddrs[18], netMsg.ToString()); err == nil {
		t.Fatal("success produce message with invalid method")
	}
	invalidMsg := random.NewRandom().GetString(8192 << 1)
	if err := testCustomProduceMessage(ctx, http.MethodPost, testutils.TgAddrs[18], invalidMsg); err == nil {
		t.Fatal("success produce message with invalid message")
	}

	msg, err := adapter1.Consume(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(msg.GetPayload().GetBody(), msgBytes) {
		t.Fatal("get invalid message bytes")
	}

	if err := adapter1.Produce(ctx, msg); err != nil {
		t.Fatal(err)
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
	defer func() { _ = rsp.Body.Close() }()
	if rsp.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}
	return nil
}

func TestSettings(t *testing.T) {
	t.Parallel()

	sett := NewSettings(nil)
	if sett.GetAdapterSettings() == nil {
		t.Fatal("invalid adapter settings")
	}
}
