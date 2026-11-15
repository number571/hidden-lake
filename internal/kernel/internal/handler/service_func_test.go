// nolint: err113
package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/number571/hidden-lake/internal/kernel/pkg/app/config"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/closer"
	"github.com/number571/hidden-lake/pkg/network/adapters/tcp"
	"github.com/number571/hidden-lake/pkg/network/handler"
	"github.com/number571/hidden-lake/pkg/network/request"
	"github.com/number571/hidden-lake/pkg/network/response"
	testutils "github.com/number571/hidden-lake/test/utils"

	anonymity "github.com/number571/go-peer/pkg/anonymity/qb"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/logger"
)

func TestHandleServiceFunc(t *testing.T) {
	t.Parallel()

	rspMsg := "hello, world!"

	mux := http.NewServeMux()
	mux.HandleFunc("/rsp-mode-on", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(hlk_settings.CHeaderResponseMode, hlk_settings.CHeaderResponseModeON)
		_, _ = fmt.Fprint(w, rspMsg)
	})
	mux.HandleFunc("/rsp-mode-off", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(hlk_settings.CHeaderResponseMode, hlk_settings.CHeaderResponseModeOFF)
		_, _ = fmt.Fprint(w, rspMsg)
	})
	mux.HandleFunc("/rsp-mode-unknown", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(hlk_settings.CHeaderResponseMode, "unknown")
		_, _ = fmt.Fprint(w, rspMsg)
	})

	addr := testutils.TgAddrs[49]
	srv := &http.Server{
		Addr:        addr,
		Handler:     mux,
		ReadTimeout: time.Second,
	}
	defer func() { _ = srv.Close() }()
	go func() { _ = srv.ListenAndServe() }()

	time.Sleep(200 * time.Millisecond)

	logger := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	ctx := context.Background()
	cfg := &tsConfig{fServiceAddr: addr}
	pubKey := tgPrivKey2.GetPubKey()
	handler := HandleServiceFunc(cfg, logger)

	reqx := request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost("hidden-some-host-not-found").
		WithPath("/rsp-mode-on").
		Build()

	if _, err := handler(ctx, pubKey, reqx); err == nil {
		t.Fatal("success handle request with invalid service")
	}

	reqy := request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost("hidden-some-host-failed").
		WithPath("/rsp-mode-on").
		Build()

	if _, err := handler(ctx, pubKey, reqy); err == nil {
		t.Fatal("success handle request with invalid do request")
	}

	req := request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost("hidden-some-host-ok").
		WithPath("/rsp-mode-on").
		Build()

	rsp, err := handler(ctx, pubKey, req)
	if err != nil {
		t.Fatal(err)
	}

	if string(rsp.GetBody()) != rspMsg {
		t.Fatal("string(rsp.GetBody()) != rspMsg")
	}

	req2 := request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost("hidden-some-host-ok").
		WithPath("/rsp-mode-off").
		Build()

	rsp2Bytes, err := handler(ctx, pubKey, req2)
	if err != nil {
		t.Fatal(err)
	}
	if rsp2Bytes != nil {
		t.Fatal("rsp2Bytes != nil")
	}

	req3 := request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost("hidden-some-host-ok").
		WithPath("/rsp-mode-unknown").
		Build()

	if _, err := handler(ctx, pubKey, req3); err == nil {
		t.Fatal("success response with unknown response mode")
	}
}

func testCleanHLS() {
	_ = os.RemoveAll(fmt.Sprintf(tcPathConfigTemplate, 9))
	for i := 0; i < 2; i++ {
		_ = os.RemoveAll(fmt.Sprintf(tcPathDBTemplate, 9+i))
	}
}

// client -> HLK -> server --\
// client <- HLK <- server <-/
func TestHLS(t *testing.T) {
	t.Parallel()

	testCleanHLS()
	defer testCleanHLS()

	// server
	srv := testStartServerHTTP(testutils.TgAddrs[5])
	defer func() { _ = srv.Close() }()

	// service
	nodeService, nodeCancel, err := testStartNodeHLS()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		nodeCancel()
		closer := closer.NewCloser(nodeService.GetKVDatabase())
		_ = closer.Close()
	}()

	// client
	nodeClient, clientCancel, err := testStartClientHLS()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		clientCancel()
		closer := closer.NewCloser(nodeClient.GetKVDatabase())
		_ = closer.Close()
	}()
}

// HLK

func testStartNodeHLS() (anonymity.INode, context.CancelFunc, error) {
	rawCFG := &config.SConfig{
		FServices: map[string]string{
			tcServiceAddressInHLS: testutils.TgAddrs[5],
		},
		FSettings: &config.SConfigSettings{
			FMessageSizeBytes: tcMessageSize,
			FWorkSizeBits:     tcWorkSize,
			FQueuePeriodMS:    tcQueuePeriod,
			FFetchTimeoutMS:   tcFetchTimeout,
		},
	}

	cfg, err := config.BuildConfig(fmt.Sprintf(tcPathConfigTemplate, 9), rawCFG)
	if err != nil {
		return nil, nil, err
	}

	logger := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	node, ctx, cancel := testRunNewNode(fmt.Sprintf(tcPathDBTemplate, 9), testutils.TgAddrs[4], handler.RequestHandler(HandleServiceFunc(cfg, logger)))
	if node == nil {
		return nil, nil, errors.New("node is not running")
	}

	node.GetKeysContainer().Add(tgPrivKey1.GetPubKey())

	networkNode := node.GetAdapter().(tcp.ITCPAdapter).GetConnKeeper().GetNetworkNode()
	go func() { _ = networkNode.Run(ctx) }()

	return node, cancel, nil
}

// CLIENT

func testStartClientHLS() (anonymity.INode, context.CancelFunc, error) {
	time.Sleep(time.Second)

	node, ctx, cancel := testRunNewNode(fmt.Sprintf(tcPathDBTemplate, 10), "", func(ctx context.Context, i anonymity.INode, ik layer2.IParticipantKey, b []byte) ([]byte, error) {
		return nil, nil
	})
	if node == nil {
		return nil, cancel, errors.New("node is not running")
	}
	node.GetKeysContainer().Add(tgPrivKey1.GetPubKey())

	networkNode := node.GetAdapter().(tcp.ITCPAdapter).GetConnKeeper().GetNetworkNode()
	if err := networkNode.AddConnection(ctx, testutils.TgAddrs[4]); err != nil {
		return nil, cancel, err
	}

	pld := request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost(tcServiceAddressInHLS).
		WithPath("/echo").
		WithHead(map[string]string{
			"Content-Type": "application/json",
		}).
		WithBody([]byte(`{"message": "hello, world!"}`)).
		Build().
		ToBytes()

	respBytes, err := node.FetchPayload(ctx, tgPrivKey1.GetPubKey(), pld)
	if err != nil {
		return node, cancel, err
	}

	resp, err := response.LoadResponse(respBytes)
	if err != nil {
		return node, cancel, err
	}

	body := resp.GetBody()
	if string(body) != "{\"echo\":\"hello, world!\",\"error\":0}\n" {
		return node, cancel, fmt.Errorf("result does not match; got '%s'", string(body))
	}

	return node, cancel, nil
}

// A sender-controlled path must not redirect the request to another host, and a
// sender-supplied Hlk-Sender-Name must not pass through.
func TestHandleServiceFuncRejectsHostInjection(t *testing.T) {
	t.Parallel()

	var attackerHit int32
	attacker := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&attackerHit, 1)
		_, _ = w.Write([]byte("pwned"))
	}))
	defer attacker.Close()
	attackerHost := strings.TrimPrefix(attacker.URL, "http://")

	log := logger.NewLogger(logger.NewSettings(&logger.SSettings{}), func(_ logger.ILogArg) string { return "" })
	handler := HandleServiceFunc(&tsConfig{fServiceAddr: "127.0.0.1:1"}, log)

	for _, path := range []string{"@" + attackerHost + "/x", "@" + attackerHost, "@evil.example/x"} {
		req := request.NewRequestBuilder().
			WithMethod(http.MethodGet).
			WithHost("hidden-some-host-ok").
			WithPath(path).
			Build()
		if _, err := handler(context.Background(), tgPrivKey2.GetPubKey(), req); !errors.Is(err, ErrInvalidServiceURL) {
			t.Fatalf("path %q: expected ErrInvalidServiceURL, got %v", path, err)
		}
	}
	if atomic.LoadInt32(&attackerHit) != 0 {
		t.Fatal("attacker host was contacted")
	}
}

func TestHandleServiceFuncOverridesForgedSenderName(t *testing.T) {
	t.Parallel()

	got := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got <- r.Header.Get(hlk_settings.CHeaderSenderName)
		w.Header().Set(hlk_settings.CHeaderResponseMode, hlk_settings.CHeaderResponseModeON)
	}))
	defer srv.Close()

	log := logger.NewLogger(logger.NewSettings(&logger.SSettings{}), func(_ logger.ILogArg) string { return "" })
	handler := HandleServiceFunc(&tsConfig{fServiceAddr: strings.TrimPrefix(srv.URL, "http://")}, log)

	// tgPrivKey2 is the friend "abc" in tsConfig; the forged header must be dropped.
	req := request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost("hidden-some-host-ok").
		WithPath("/x").
		WithHead(map[string]string{hlk_settings.CHeaderSenderName: "forged-alias"}).
		Build()
	if _, err := handler(context.Background(), tgPrivKey2.GetPubKey(), req); err != nil {
		t.Fatal(err)
	}
	if name := <-got; name != "abc" {
		t.Fatalf("service saw sender name %q, want \"abc\"", name)
	}
}
