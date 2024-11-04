package handler

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/number571/hidden-lake/internal/applications/remoter/pkg/client"
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestIncomingExecHTTP(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, service := testRunService(testutils.TgAddrs[40])
	defer service.Close()

	testRunNewNodes(
		ctx,
		testutils.TgAddrs[41],
		testutils.TgAddrs[42],
		testutils.TgAddrs[40],
	)

	hlsClient := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			"http://"+testutils.TgAddrs[41],
			&http.Client{Timeout: time.Minute},
		),
	)

	hlrClient := client.NewClient(
		client.NewBuilder(tcPassword),
		client.NewRequester(hlsClient),
	)

	msg := "hello, world!"
	rsp, err := hlrClient.Exec(ctx, "test_recv", "echo", msg)
	if err != nil {
		t.Error(err)
		return
	}

	if strings.TrimSpace(string(rsp)) != msg {
		t.Error("get invalid response")
		return
	}
}
