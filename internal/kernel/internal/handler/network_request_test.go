package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	anonymity "github.com/number571/go-peer/pkg/anonymity/qb"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/kernel/pkg/app/config"
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	"github.com/number571/hidden-lake/internal/utils/closer"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	"github.com/number571/hidden-lake/pkg/adapters/tcp"
	"github.com/number571/hidden-lake/pkg/network/handler"
	"github.com/number571/hidden-lake/pkg/network/request"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestHandleRequestAPI2(t *testing.T) {
	t.Parallel()

	httpLogger := std_logger.NewStdLogger(
		func() std_logger.ILogging {
			logging, err := std_logger.LoadLogging([]string{})
			if err != nil {
				panic(err)
			}
			return logging
		}(),
		func(_ logger.ILogArg) string {
			return ""
		},
	)

	ctx := context.Background()

	handler := HandleNetworkRequestAPI(
		ctx,
		&tsConfig{},
		httpLogger,
		newTsHiddenLakeNode(newTsNode(true, true, true)),
	)
	if err := requestAPIRequestPutOK(handler); err != nil {
		t.Fatal(err)
	}
	if err := requestAPIRequestPostOK(handler); err != nil {
		t.Fatal(err)
	}

	if err := requestAPIRequestReqData(handler); err == nil {
		t.Fatal("request success with invalid reqData")
	}
	if err := requestAPIRequestNotFound(handler); err == nil {
		t.Fatal("request success with not found alias_name")
	}
	if err := requestAPIRequestDecode(handler); err == nil {
		t.Fatal("request success with invalid decode")
	}
	if err := requestAPIRequestMethod(handler); err == nil {
		t.Fatal("request success with invalid method")
	}

	handlerx := HandleNetworkRequestAPI(
		ctx,
		&tsConfig{},
		httpLogger,
		newTsHiddenLakeNode(newTsNode(false, false, true)),
	)
	if err := requestAPIRequestPutOK(handlerx); err == nil {
		t.Fatalf("request success with put error")
	}
	if err := requestAPIRequestPostOK(handlerx); err == nil {
		t.Fatalf("request success with post error")
	}

	handlery := HandleNetworkRequestAPI(
		ctx,
		&tsConfig{},
		httpLogger,
		newTsHiddenLakeNode(newTsNode(true, true, false)),
	)
	if err := requestAPIRequestPostOK(handlery); err == nil {
		t.Fatal("request success with post error (load response)")
	}
}

func requestAPIRequestPutOK(handler http.HandlerFunc) error {
	request := &request.SRequest{}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/?friend=abc", bytes.NewBuffer(encoding.SerializeJSON(request)))

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func requestAPIRequestPostOK(handler http.HandlerFunc) error {
	request := &request.SRequest{}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/?friend=abc", bytes.NewBuffer(encoding.SerializeJSON(request)))

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func requestAPIRequestReqData(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/?friend=abc", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func requestAPIRequestNotFound(handler http.HandlerFunc) error {
	request := &request.SRequest{}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/?friend=notfound", bytes.NewBuffer(encoding.SerializeJSON(request)))

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func requestAPIRequestDecode(handler http.HandlerFunc) error {
	request := &request.SRequest{}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/?friend=abc", bytes.NewBuffer(bytes.Join(
		[][]byte{
			[]byte{1},
			encoding.SerializeJSON(request),
		},
		[]byte{},
	)))

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func requestAPIRequestMethod(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func TestHandleRequestAPI(t *testing.T) {
	t.Parallel()

	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 7)
	pathDB := fmt.Sprintf(tcPathDBTemplate, 7)

	_, node, ctx, cancel, srv := testAllCreate(pathCfg, pathDB, testutils.TgAddrs[9])
	defer testAllFree(node, cancel, srv, pathCfg, pathDB)

	pushNode, pushCancel, pushSrv := testAllPushCreate(pathCfg, pathDB)
	defer testAllPushFree(pushNode, pushCancel, pushSrv, pathCfg, pathDB)

	client := hlk_client.NewClient(
		hlk_client.NewBuilder(),
		hlk_client.NewRequester(
			testutils.TgAddrs[9],
			&http.Client{Timeout: time.Minute},
		),
	)

	networkNode := node.GetAdapter().(tcp.ITCPAdapter).GetConnKeeper().GetNetworkNode()
	_ = networkNode.AddConnection(ctx, testutils.TgAddrs[11])
	node.GetMapPubKeys().SetPubKey(tgPrivKey1.GetPubKey())

	testSend(t, client)
	testFetch(t, client)
}

func testSend(t *testing.T, client hlk_client.IClient) {
	err := client.SendRequest(
		context.Background(),
		"test_recvr",
		request.NewRequestBuilder().
			WithMethod(http.MethodGet).
			WithHost(tcServiceAddressInHLS).
			WithPath("/echo").
			WithHead(map[string]string{
				"Content-Type": "application/json",
			}).
			WithBody([]byte(`{"message": "hello, world!"}`)).
			Build(),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func testFetch(t *testing.T, client hlk_client.IClient) {
	res, err := client.FetchRequest(
		context.Background(),
		"test_recvr",
		request.NewRequestBuilder().
			WithMethod(http.MethodGet).
			WithHost(tcServiceAddressInHLS).
			WithPath("/echo").
			WithHead(map[string]string{
				"Content-Type": "application/json",
			}).
			WithBody([]byte(`{"message": "hello, world!"}`)).
			Build(),
	)
	if err != nil {
		t.Fatal(err)
	}

	body := res.GetBody()
	if string(body) != "{\"echo\":\"hello, world!\",\"error\":0}\n" {
		t.Fatalf("result does not match; got '%s'", string(body))
	}
}

func testAllPushCreate(pathCfg, pathDB string) (anonymity.INode, context.CancelFunc, *http.Server) {
	_ = os.RemoveAll(pathCfg + "_push1")
	_ = os.RemoveAll(pathDB + "_push1")

	pushNode, cancel := testNewPushNode(pathCfg+"_push1", pathDB+"_push1")
	if pushNode == nil {
		return nil, cancel, nil
	}

	pushSrv := testStartServerHTTP(testutils.TgAddrs[10])
	time.Sleep(200 * time.Millisecond)
	return pushNode, cancel, pushSrv
}

func testAllPushFree(node anonymity.INode, cancel context.CancelFunc, srv *http.Server, pathCfg, pathDB string) {
	defer func() {
		_ = os.RemoveAll(pathCfg + "_push1")
		_ = os.RemoveAll(pathDB + "_push1")
	}()
	cancel()
	closer := closer.NewCloser(
		srv,
		node.GetKVDatabase(),
	)
	_ = closer.Close()
}

func testNewPushNode(cfgPath, dbPath string) (anonymity.INode, context.CancelFunc) {
	node, ctx, cancel := testRunNewNode(dbPath, testutils.TgAddrs[11])
	rawCFG := &config.SConfig{
		FSettings: &config.SConfigSettings{
			FMessageSizeBytes: tcMessageSize,
			FWorkSizeBits:     tcWorkSize,
			FQueuePeriodMS:    tcQueuePeriod,
			FFetchTimeoutMS:   tcFetchTimeout,
		},
		FServices: map[string]string{
			tcServiceAddressInHLS: testutils.TgAddrs[10],
		},
	}

	cfg, err := config.BuildConfig(cfgPath, rawCFG)
	if err != nil {
		return nil, cancel
	}

	logger := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	node.HandleFunc(
		build.GetSettings().FProtoMask.FService,
		handler.RequestHandler(HandleServiceFunc(cfg, logger)),
	)
	node.GetMapPubKeys().SetPubKey(tgPrivKey1.GetPubKey())

	networkNode := node.GetAdapter().(tcp.ITCPAdapter).GetConnKeeper().GetNetworkNode()
	go func() { _ = networkNode.Run(ctx) }()

	return node, cancel
}
