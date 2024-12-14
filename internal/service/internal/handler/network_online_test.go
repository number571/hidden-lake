package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/number571/go-peer/pkg/logger"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	"github.com/number571/hidden-lake/pkg/adapters/http/client"
)

func TestHandleOnlineAPI2(t *testing.T) {
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

	epClients := []client.IClient{
		client.NewClient(&tsRequester{}),
	}

	handler := HandleNetworkOnlineAPI(ctx, httpLogger, epClients)
	if err := onlineAPIRequestOK(handler); err != nil {
		t.Error(err)
		return
	}
	if err := onlineAPIRequestDeleteOK(handler); err != nil {
		t.Error(err)
		return
	}

	if err := onlineAPIRequestMethod(handler); err == nil {
		t.Error("request success with invalid method")
		return
	}

	epClientsx := []client.IClient{
		client.NewClient(&tsRequester{fWithFail: true}),
	}

	handlerx := HandleNetworkOnlineAPI(ctx, httpLogger, epClientsx)
	if err := onlineAPIRequestDeleteOK(handlerx); err == nil {
		t.Error("request success with delete error")
		return
	}
}

func onlineAPIRequestOK(handler http.HandlerFunc) error {
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

func onlineAPIRequestDeleteOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", strings.NewReader("127.0.0.1:9999"))

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

func onlineAPIRequestMethod(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)

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

// func TestHandleOnlineAPI(t *testing.T) {
// 	t.Parallel()

// 	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 6)
// 	pathDB := fmt.Sprintf(tcPathDBTemplate, 6)

// 	_, node, ctx, cancel, srv := testAllCreate(pathCfg, pathDB, testutils.TgAddrs[12])
// 	defer testAllFree(node, cancel, srv, pathCfg, pathDB)

// 	pushNode, pushCancel := testAllOnlineCreate(pathCfg, pathDB)
// 	defer testAllOnlineFree(pushNode, pushCancel, pathCfg, pathDB)

// 	client := hls_client.NewClient(
// 		hls_client.NewBuilder(),
// 		hls_client.NewRequester(
// 			"http://"+testutils.TgAddrs[12],
// 			&http.Client{Timeout: time.Minute},
// 		),
// 	)

// 	networkNode := node.GetAdapter().(tcp.ITCPAdapter).GetConnKeeper().GetNetworkNode()
// 	_ = networkNode.AddConnection(ctx, testutils.TgAddrs[13])
// 	node.GetMapPubKeys().SetPubKey(tgPrivKey1.GetPubKey())

// 	testGetOnlines(t, client, node)
// 	testDelOnline(t, client, testutils.TgAddrs[13])
// }

// func testGetOnlines(t *testing.T, client hls_client.IClient, node anonymity.INode) {
// 	onlines, err := client.GetOnlines(context.Background())
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	if len(onlines) != 1 {
// 		t.Error("length of onlines != 1")
// 		return
// 	}

// 	networkNode := node.GetAdapter().(tcp.ITCPAdapter).GetConnKeeper().GetNetworkNode()
// 	if _, ok := networkNode.GetConnections()[onlines[0]]; !ok {
// 		t.Error("online address is invalid")
// 		return
// 	}
// }

// func testDelOnline(t *testing.T, client hls_client.IClient, addr string) {
// 	err := client.DelOnline(context.Background(), addr)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	onlines, err := client.GetOnlines(context.Background())
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	if len(onlines) != 0 {
// 		t.Error("length of onlines != 0")
// 		return
// 	}
// }

// func testAllOnlineCreate(pathCfg, pathDB string) (anonymity.INode, context.CancelFunc) {
// 	os.RemoveAll(pathCfg + "_push2")
// 	os.RemoveAll(pathDB + "_push2")

// 	pushNode, cancel := testOnlinePushNode(pathCfg+"_push2", pathDB+"_push2")
// 	if pushNode == nil {
// 		return nil, nil
// 	}

// 	time.Sleep(200 * time.Millisecond)
// 	return pushNode, cancel
// }

// func testAllOnlineFree(node anonymity.INode, cancel context.CancelFunc, pathCfg, pathDB string) {
// 	defer func() {
// 		os.RemoveAll(pathCfg + "_push2")
// 		os.RemoveAll(pathDB + "_push2")
// 	}()
// 	cancel()
// 	_ = closer.CloseAll([]io.Closer{
// 		node.GetKVDatabase(),
// 	})
// }

// func testOnlinePushNode(cfgPath, dbPath string) (anonymity.INode, context.CancelFunc) {
// 	node, ctx, cancel := testRunNewNode(dbPath, testutils.TgAddrs[13])

// 	cfg, err := config.BuildConfig(cfgPath, &config.SConfig{
// 		FSettings: &config.SConfigSettings{
// 			FMessageSizeBytes: tcMessageSize,
// 			FWorkSizeBits:     tcWorkSize,
// 			FQueuePeriodMS:    tcQueuePeriod,
// 			FFetchTimeoutMS:   tcFetchTimeout,
// 		},
// 	})
// 	if err != nil {
// 		return nil, cancel
// 	}

// 	logger := logger.NewLogger(
// 		logger.NewSettings(&logger.SSettings{}),
// 		func(_ logger.ILogArg) string { return "" },
// 	)

// 	node.HandleFunc(
// 		build.GSettings.FProtoMask.FService,
// 		handler.RequestHandler(HandleServiceFunc(cfg, logger)),
// 	)
// 	node.GetMapPubKeys().SetPubKey(tgPrivKey1.GetPubKey())

// 	networkNode := node.GetAdapter().(tcp.ITCPAdapter).GetConnKeeper().GetNetworkNode()
// 	go func() { _ = networkNode.Run(ctx) }()
// 	return node, cancel
// }
