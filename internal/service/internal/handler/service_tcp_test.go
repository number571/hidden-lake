// nolint: goerr113
package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/hidden-lake/internal/service/internal/config"
	"github.com/number571/hidden-lake/internal/service/pkg/request"
	"github.com/number571/hidden-lake/internal/service/pkg/response"
	pkg_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/closer"
	testutils "github.com/number571/hidden-lake/test/utils"

	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/types"
)

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
		_ = closer.CloseAll([]types.ICloser{
			nodeService.GetKVDatabase(),
			nodeService.GetNetworkNode(),
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
		_ = closer.CloseAll([]types.ICloser{
			nodeClient.GetKVDatabase(),
			nodeClient.GetNetworkNode(),
		})
	}()
}

// HLS

func testStartNodeHLS() (anonymity.INode, context.CancelFunc, error) {
	rawCFG := &config.SConfig{
		FServices: map[string]*config.SService{
			tcServiceAddressInHLS: {FHost: testutils.TgAddrs[5]},
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

	node.HandleFunc(
		pkg_settings.CServiceMask,
		HandleServiceTCP(cfg),
	)
	node.GetMapPubKeys().SetPubKey(tgPrivKey1.GetPubKey())

	go func() {
		_ = node.GetNetworkNode().Listen(ctx)
	}()

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

	if err := node.GetNetworkNode().AddConnection(ctx, testutils.TgAddrs[4]); err != nil {
		return nil, cancel, err
	}

	pld := payload.NewPayload32(
		pkg_settings.CServiceMask,
		request.NewRequest(http.MethodGet, tcServiceAddressInHLS, "/echo").
			WithHead(map[string]string{
				"Content-Type": "application/json",
			}).
			WithBody([]byte(`{"message": "hello, world!"}`)).
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
