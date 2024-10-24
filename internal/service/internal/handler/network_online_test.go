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
	pkg_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/closer"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestHandleOnlineAPI(t *testing.T) {
	t.Parallel()

	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 6)
	pathDB := fmt.Sprintf(tcPathDBTemplate, 6)

	_, node, ctx, cancel, srv := testAllCreate(pathCfg, pathDB, testutils.TgAddrs[12])
	defer testAllFree(node, cancel, srv, pathCfg, pathDB)

	pushNode, pushCancel := testAllOnlineCreate(pathCfg, pathDB)
	defer testAllOnlineFree(pushNode, pushCancel, pathCfg, pathDB)

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			"http://"+testutils.TgAddrs[12],
			&http.Client{Timeout: time.Minute},
		),
	)

	_ = node.GetNetworkNode().AddConnection(ctx, testutils.TgAddrs[13])
	node.GetMapPubKeys().SetPubKey(tgPrivKey1.GetPubKey())

	testGetOnlines(t, client, node)
	testDelOnline(t, client, testutils.TgAddrs[13])
}

func testGetOnlines(t *testing.T, client hls_client.IClient, node anonymity.INode) {
	onlines, err := client.GetOnlines(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if len(onlines) != 1 {
		t.Error("length of onlines != 1")
		return
	}

	if _, ok := node.GetNetworkNode().GetConnections()[onlines[0]]; !ok {
		t.Error("online address is invalid")
		return
	}
}

func testDelOnline(t *testing.T, client hls_client.IClient, addr string) {
	err := client.DelOnline(context.Background(), addr)
	if err != nil {
		t.Error(err)
		return
	}

	onlines, err := client.GetOnlines(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if len(onlines) != 0 {
		t.Error("length of onlines != 0")
		return
	}
}

func testAllOnlineCreate(pathCfg, pathDB string) (anonymity.INode, context.CancelFunc) {
	os.RemoveAll(pathCfg + "_push2")
	os.RemoveAll(pathDB + "_push2")

	pushNode, cancel := testOnlinePushNode(pathCfg+"_push2", pathDB+"_push2")
	if pushNode == nil {
		return nil, nil
	}

	time.Sleep(200 * time.Millisecond)
	return pushNode, cancel
}

func testAllOnlineFree(node anonymity.INode, cancel context.CancelFunc, pathCfg, pathDB string) {
	defer func() {
		os.RemoveAll(pathCfg + "_push2")
		os.RemoveAll(pathDB + "_push2")
	}()
	cancel()
	_ = closer.CloseAll([]types.ICloser{
		node.GetKVDatabase(),
		node.GetNetworkNode(),
	})
}

func testOnlinePushNode(cfgPath, dbPath string) (anonymity.INode, context.CancelFunc) {
	node, ctx, cancel := testRunNewNode(dbPath, testutils.TgAddrs[13])

	cfg, err := config.BuildConfig(cfgPath, &config.SConfig{
		FSettings: &config.SConfigSettings{
			FMessageSizeBytes: tcMessageSize,
			FWorkSizeBits:     tcWorkSize,
			FQueuePeriodMS:    tcQueuePeriod,
			FFetchTimeoutMS:   tcFetchTimeout,
		},
	})
	if err != nil {
		return nil, cancel
	}

	node.HandleFunc(
		pkg_settings.CServiceMask,
		HandleServiceTCP(cfg),
	)
	node.GetMapPubKeys().SetPubKey(tgPrivKey1.GetPubKey())

	go func() { _ = node.GetNetworkNode().Listen(ctx) }()

	return node, cancel
}
