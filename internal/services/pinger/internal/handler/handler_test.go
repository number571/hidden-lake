package handler

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/build"
	hls_settings "github.com/number571/hidden-lake/internal/services/pinger/pkg/settings"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	hls_client "github.com/number571/hidden-lake/pkg/api/services/pinger/client"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestHandler(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
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

	server := testInitInternalServiceHTTP(ctx, httpLogger, newTsHLKClient(0), testutils.TgAddrs[13])
	defer func() { _ = server.Close() }()

	go func() { _ = server.ListenAndServe() }()
	time.Sleep(100 * time.Millisecond)

	client := hls_client.NewClient(
		hls_client.NewRequester(
			testutils.TgAddrs[13],
			&http.Client{Timeout: time.Second},
		),
	)

	if err := client.GetIndex(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := client.PingFriend(context.Background(), "abc"); err != nil {
		t.Fatal(err)
	}
}

func testInitInternalServiceHTTP(
	pCtx context.Context,
	pLogger logger.ILogger,
	pHlkClient hlk_client.IClient,
	pAddress string,
) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc(
		hls_settings.CHandleIndexPath,
		HandleIndexAPI(pLogger),
	) // GET

	mux.HandleFunc(
		hls_settings.CHandleCommandPingPath,
		HandleCommandPingAPI(pCtx, pLogger, pHlkClient),
	) // GET

	buildSettings := build.GetSettings()
	return &http.Server{ // nolint: gosec
		Addr:         pAddress,
		Handler:      http.TimeoutHandler(mux, buildSettings.GetHttpCallbackTimeout(), "handle timeout"),
		ReadTimeout:  buildSettings.GetHttpCallbackTimeout(),
		WriteTimeout: buildSettings.GetHttpCallbackTimeout(),
	}
}
