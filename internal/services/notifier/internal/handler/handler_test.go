package handler

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/services/notifier/internal/message"
	hls_settings "github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/broker"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	hls_client "github.com/number571/hidden-lake/pkg/api/services/notifier/client"
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

	msgBroker := broker.NewDataBroker(256, 256)
	server := testInitInternalServiceHTTP(
		ctx,
		httpLogger,
		newTsHLKClient(true, true, false),
		msgBroker,
		testutils.TgAddrs[24],
	)
	defer func() { _ = server.Close() }()

	go func() { _ = server.ListenAndServe() }()
	time.Sleep(100 * time.Millisecond)

	client := hls_client.NewClient(
		testutils.TgAddrs[24],
		&http.Client{Timeout: time.Second},
	)

	if err := client.GetIndex(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := client.SendNotification(context.Background(), "abc", []byte("hello, world!")); err != nil {
		t.Fatal(err)
	}

	msg := []byte("hello, world!")

	go func() {
		time.Sleep(100 * time.Millisecond)
		msgBroker.Produce(message.NewMessageContainer("abc", msg))
	}()
	_, gotMsg, err := client.RecvNotification(context.Background(), "sid")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(gotMsg, msg) {
		t.Fatal("gotMsg.ToString() != msg.ToString()")
	}
}

func testInitInternalServiceHTTP(
	pCtx context.Context,
	pLogger logger.ILogger,
	pHlkClient hlk_client.IClient,
	pMsgBroker broker.IDataBroker,
	pAddress string,
) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc(
		hls_settings.CHandleIndexPath,
		HandleIndexAPI(pLogger),
	) // GET

	mux.HandleFunc(
		hls_settings.CHandleChanNotifyPath,
		HandleChanNotifyAPI(pCtx, pLogger, pHlkClient),
	) // POST

	mux.HandleFunc(
		hls_settings.CHandleChanListenPath,
		HandleChanListenAPI(pCtx, pLogger, pMsgBroker),
	) // GET

	buildSettings := build.GetSettings()
	return &http.Server{ // nolint: gosec
		Addr:         pAddress,
		Handler:      http.TimeoutHandler(mux, buildSettings.GetHttpCallbackTimeout(), "handle timeout"),
		ReadTimeout:  buildSettings.GetHttpCallbackTimeout(),
		WriteTimeout: buildSettings.GetHttpCallbackTimeout(),
	}
}
