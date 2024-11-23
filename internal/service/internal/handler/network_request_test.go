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

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/service/pkg/app/config"
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
	pkg_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/closer"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	"github.com/number571/hidden-lake/pkg/handler"
	"github.com/number571/hidden-lake/pkg/request"
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
		newTsHiddenLakeNode(newTsNode(true, true, true, true)),
	)
	if err := requestAPIRequestPutOK(handler); err != nil {
		t.Error(err)
		return
	}
	if err := requestAPIRequestPostOK(handler); err != nil {
		t.Error(err)
		return
	}

	if err := requestAPIRequestReqData(handler); err == nil {
		t.Error("request success with invalid reqData")
		return
	}
	if err := requestAPIRequestNotFound(handler); err == nil {
		t.Error("request success with not found alias_name")
		return
	}
	if err := requestAPIRequestDecode(handler); err == nil {
		t.Error("request success with invalid decode")
		return
	}
	if err := requestAPIRequestMethod(handler); err == nil {
		t.Error("request success with invalid method")
		return
	}

	handlerx := HandleNetworkRequestAPI(
		ctx,
		&tsConfig{},
		httpLogger,
		newTsHiddenLakeNode(newTsNode(true, false, false, true)),
	)
	if err := requestAPIRequestPutOK(handlerx); err == nil {
		t.Error("request success with put error")
		return
	}
	if err := requestAPIRequestPostOK(handlerx); err == nil {
		t.Error("request success with post error")
		return
	}

	handlery := HandleNetworkRequestAPI(
		ctx,
		&tsConfig{},
		httpLogger,
		newTsHiddenLakeNode(newTsNode(true, true, true, false)),
	)
	if err := requestAPIRequestPostOK(handlery); err == nil {
		t.Error("request success with post error (load response)")
		return
	}
}

func requestAPIRequestPutOK(handler http.HandlerFunc) error {
	request := pkg_settings.SRequest{
		FReceiver: "abc",
		FReqData:  &request.SRequest{},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(encoding.SerializeJSON(request)))

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func requestAPIRequestPostOK(handler http.HandlerFunc) error {
	request := pkg_settings.SRequest{
		FReceiver: "abc",
		FReqData:  &request.SRequest{},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(encoding.SerializeJSON(request)))

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func requestAPIRequestReqData(handler http.HandlerFunc) error {
	request := pkg_settings.SRequest{
		FReceiver: "abc",
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(encoding.SerializeJSON(request)))

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func requestAPIRequestNotFound(handler http.HandlerFunc) error {
	request := pkg_settings.SRequest{
		FReceiver: "notfound",
		FReqData:  &request.SRequest{},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(encoding.SerializeJSON(request)))

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func requestAPIRequestDecode(handler http.HandlerFunc) error {
	request := pkg_settings.SRequest{
		FReceiver: "abc",
		FReqData:  &request.SRequest{},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bytes.Join(
		[][]byte{
			[]byte{1},
			encoding.SerializeJSON(request),
		},
		[]byte{},
	)))

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

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
	defer res.Body.Close()

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

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			"http://"+testutils.TgAddrs[9],
			&http.Client{Timeout: time.Minute},
		),
	)

	_ = node.GetNetworkNode().AddConnection(ctx, testutils.TgAddrs[11])
	node.GetMapPubKeys().SetPubKey(tgPrivKey1.GetPubKey())

	testSend(t, client)
	testFetch(t, client)
}

func testSend(t *testing.T, client hls_client.IClient) {
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
		t.Error(err)
		return
	}
}

func testFetch(t *testing.T, client hls_client.IClient) {
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
		t.Error(err)
		return
	}

	body := res.GetBody()
	if string(body) != "{\"echo\":\"hello, world!\",\"error\":0}\n" {
		t.Errorf("result does not match; got '%s'", string(body))
		return
	}
}

func testAllPushCreate(pathCfg, pathDB string) (anonymity.INode, context.CancelFunc, *http.Server) {
	os.RemoveAll(pathCfg + "_push1")
	os.RemoveAll(pathDB + "_push1")

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
		os.RemoveAll(pathCfg + "_push1")
		os.RemoveAll(pathDB + "_push1")
	}()
	cancel()
	_ = closer.CloseAll([]types.ICloser{
		srv,
		node.GetKVDatabase(),
		node.GetNetworkNode(),
	})
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
		build.GSettings.FProtoMask.FService,
		handler.RequestHandler(HandleServiceTCP(cfg, logger)),
	)
	node.GetMapPubKeys().SetPubKey(tgPrivKey1.GetPubKey())

	go func() { _ = node.GetNetworkNode().Listen(ctx) }()

	return node, cancel
}
