package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/internal/service/internal/config"
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
	"github.com/number571/hidden-lake/internal/service/pkg/request"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/closer"
	testutils "github.com/number571/hidden-lake/test/utils"
)

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

	testBroadcast(t, client)
	testFetch(t, client)
}

func testBroadcast(t *testing.T, client hls_client.IClient) {
	err := client.BroadcastRequest(
		context.Background(),
		"test_recvr",
		request.NewRequest(http.MethodGet, tcServiceAddressInHLS, "/echo").
			WithHead(map[string]string{
				"Content-Type": "application/json",
			}).
			WithBody([]byte(`{"message": "hello, world!"}`)),
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
		request.NewRequest(http.MethodGet, tcServiceAddressInHLS, "/echo").
			WithHead(map[string]string{
				"Content-Type": "application/json",
			}).
			WithBody([]byte(`{"message": "hello, world!"}`)),
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
		FServices: map[string]*config.SService{
			tcServiceAddressInHLS: {FHost: testutils.TgAddrs[10]},
		},
	}

	cfg, err := config.BuildConfig(cfgPath, rawCFG)
	if err != nil {
		return nil, cancel
	}

	node.HandleFunc(
		hls_settings.CServiceMask,
		HandleServiceTCP(cfg),
	)
	node.GetMapPubKeys().SetPubKey(tgPrivKey1.GetPubKey())

	go func() { _ = node.GetNetworkNode().Listen(ctx) }()

	return node, cancel
}
