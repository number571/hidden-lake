// nolint: goerr113
package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/service/pkg/app/config"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/closer"
	"github.com/number571/hidden-lake/pkg/adapters/tcp"
	"github.com/number571/hidden-lake/pkg/handler"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"
	testutils "github.com/number571/hidden-lake/test/utils"

	"github.com/number571/go-peer/pkg/anonymity"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/payload"
)

func TestHandleServiceFunc(t *testing.T) {
	t.Parallel()

	rspMsg := "hello, world!"

	mux := http.NewServeMux()
	mux.HandleFunc("/rsp-mode-on", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(hls_settings.CHeaderResponseMode, hls_settings.CHeaderResponseModeON)
		fmt.Fprint(w, rspMsg)
	})
	mux.HandleFunc("/rsp-mode-off", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(hls_settings.CHeaderResponseMode, hls_settings.CHeaderResponseModeOFF)
		fmt.Fprint(w, rspMsg)
	})
	mux.HandleFunc("/rsp-mode-unknown", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(hls_settings.CHeaderResponseMode, "unknown")
		fmt.Fprint(w, rspMsg)
	})

	addr := testutils.TgAddrs[49]
	srv := &http.Server{
		Addr:        addr,
		Handler:     mux,
		ReadTimeout: time.Second,
	}
	defer srv.Close()
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
		t.Error("success handle request with invalid service")
		return
	}

	reqy := request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost("hidden-some-host-failed").
		WithPath("/rsp-mode-on").
		Build()

	if _, err := handler(ctx, pubKey, reqy); err == nil {
		t.Error("success handle request with invalid do request")
		return
	}

	req := request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost("hidden-some-host-ok").
		WithPath("/rsp-mode-on").
		Build()

	rsp, err := handler(ctx, pubKey, req)
	if err != nil {
		t.Error(err)
		return
	}

	if string(rsp.GetBody()) != rspMsg {
		t.Error("string(rsp.GetBody()) != rspMsg")
		return
	}

	req2 := request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost("hidden-some-host-ok").
		WithPath("/rsp-mode-off").
		Build()

	rsp2Bytes, err := handler(ctx, pubKey, req2)
	if err != nil {
		t.Error(err)
		return
	}
	if rsp2Bytes != nil {
		t.Error("rsp2Bytes != nil")
		return
	}

	req3 := request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost("hidden-some-host-ok").
		WithPath("/rsp-mode-unknown").
		Build()

	if _, err := handler(ctx, pubKey, req3); err == nil {
		t.Error("success response with unknown response mode")
		return
	}
}

func testCleanHLS() {
	os.RemoveAll(fmt.Sprintf(tcPathConfigTemplate, 9))
	for i := 0; i < 2; i++ {
		os.RemoveAll(fmt.Sprintf(tcPathDBTemplate, 9+i))
	}
}

// client -> HLS -> server --\
// client <- HLS <- server <-/
func TestHLS(t *testing.T) {
	t.Parallel()

	testCleanHLS()
	defer testCleanHLS()

	// server
	srv := testStartServerHTTP(testutils.TgAddrs[5])
	defer srv.Close()

	// service
	nodeService, nodeCancel, err := testStartNodeHLS()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		nodeCancel()
		_ = closer.CloseAll([]io.Closer{
			nodeService.GetKVDatabase(),
		})
	}()

	// client
	nodeClient, clientCancel, err := testStartClientHLS()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		clientCancel()
		_ = closer.CloseAll([]io.Closer{
			nodeClient.GetKVDatabase(),
		})
	}()
}

// HLS

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

	node, ctx, cancel := testRunNewNode(fmt.Sprintf(tcPathDBTemplate, 9), testutils.TgAddrs[4])
	if node == nil {
		return nil, nil, errors.New("node is not running")
	}

	logger := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	node.HandleFunc(
		build.GSettings.FProtoMask.FService,
		handler.RequestHandler(HandleServiceFunc(cfg, logger)),
	)
	node.GetMapPubKeys().SetPubKey(tgPrivKey1.GetPubKey())

	networkNode := node.GetAdapter().(tcp.ITCPAdapter).GetConnKeeper().GetNetworkNode()
	go func() { _ = networkNode.Run(ctx) }()

	return node, cancel, nil
}

// CLIENT

func testStartClientHLS() (anonymity.INode, context.CancelFunc, error) {
	time.Sleep(time.Second)

	node, ctx, cancel := testRunNewNode(fmt.Sprintf(tcPathDBTemplate, 10), "")
	if node == nil {
		return nil, cancel, errors.New("node is not running")
	}
	node.GetMapPubKeys().SetPubKey(tgPrivKey1.GetPubKey())

	networkNode := node.GetAdapter().(tcp.ITCPAdapter).GetConnKeeper().GetNetworkNode()
	if err := networkNode.AddConnection(ctx, testutils.TgAddrs[4]); err != nil {
		return nil, cancel, err
	}

	pld := payload.NewPayload32(
		build.GSettings.FProtoMask.FService,
		request.NewRequestBuilder().
			WithMethod(http.MethodGet).
			WithHost(tcServiceAddressInHLS).
			WithPath("/echo").
			WithHead(map[string]string{
				"Content-Type": "application/json",
			}).
			WithBody([]byte(`{"message": "hello, world!"}`)).
			Build().
			ToBytes(),
	)

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
