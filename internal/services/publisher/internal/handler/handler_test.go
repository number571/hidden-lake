package handler

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/services/publisher/internal/database"
	"github.com/number571/hidden-lake/internal/services/publisher/internal/handler/process"
	"github.com/number571/hidden-lake/internal/services/publisher/internal/message"
	"github.com/number571/hidden-lake/internal/services/publisher/pkg/app/config"
	hls_settings "github.com/number571/hidden-lake/internal/services/publisher/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/broker"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	hls_client "github.com/number571/hidden-lake/pkg/api/services/publisher/client"
	"github.com/number571/hidden-lake/pkg/api/services/publisher/client/dto"
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
		&tsConfig{},
		newTsHLKClient(true, true, false),
		newTsDatabase(true, true),
		msgBroker,
		testutils.TgAddrs[13],
	)
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
	if _, err := client.GetPushLimit(context.Background()); err != nil {
		t.Fatal(err)
	}

	msgUID := dto.NewMessageUID()
	if _, err := client.PushMessageToBox(context.Background(), "abc", msgUID, "hello, world!"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.LoadMessageByEnum(context.Background(), "abc", 0); err != nil {
		t.Fatal(err)
	}
	if _, err := client.LoadMessageByUID(context.Background(), "abc", msgUID); err != nil {
		t.Fatal(err)
	}

	msg := dto.NewMessageBody("hello, world!", time.Now())

	go func() {
		time.Sleep(100 * time.Millisecond)
		msgBroker.Produce(message.NewMessageContainer("abc", "", msg))
	}()
	_, gotMsg, err := client.ListenMsgbox(context.Background(), "abc", "sid")
	if err != nil {
		t.Fatal(err)
	}
	if gotMsg.ToString() != msg.ToString() {
		t.Fatal("gotMsg.ToString() != msg.ToString()")
	}
}

func testInitInternalServiceHTTP(
	pCtx context.Context,
	pLogger logger.ILogger,
	pConfig config.IConfig,
	pHlkClient hlk_client.IClient,
	pDatabase database.IKVDatabase,
	pMsgBroker broker.IDataBroker,
	pAddress string,
) *http.Server {
	mux := http.NewServeMux()
	processManager := process.NewPublishProcessManager()

	mux.HandleFunc(
		hls_settings.CHandleIndexPath,
		HandleIndexAPI(pLogger),
	) // GET

	mux.HandleFunc(
		hls_settings.CHandleMsgboxSizePath,
		HandleChanSizeAPI(pCtx, pLogger, pConfig, pDatabase),
	) // GET

	mux.HandleFunc(
		hls_settings.CHandleMsgboxLoadPath,
		HandleChanLoadAPI(pCtx, pLogger, pConfig, pHlkClient, pDatabase),
	) // GET

	mux.HandleFunc(
		hls_settings.CHandleMsgboxPushPath,
		HandleChanPushAPI(pCtx, pLogger, pConfig, pHlkClient, pDatabase, processManager),
	) // POST

	mux.HandleFunc(
		hls_settings.CHandleMsgboxListenPath,
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
