package handler

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/services/messenger/internal/database"
	"github.com/number571/hidden-lake/internal/services/messenger/internal/message"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/app/config"
	hls_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	hls_client "github.com/number571/hidden-lake/pkg/api/services/messenger/client"
	"github.com/number571/hidden-lake/pkg/api/services/messenger/client/dto"
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

	msgBroker := message.NewMessageBroker()
	server := testInitInternalServiceHTTP(
		ctx,
		httpLogger,
		&tsConfig{},
		newTsHLKClient(true, true, true),
		newTsDatabase(true, true),
		msgBroker,
		testutils.TgAddrs[24],
	)
	defer func() { _ = server.Close() }()

	go func() { _ = server.ListenAndServe() }()
	time.Sleep(100 * time.Millisecond)

	client := hls_client.NewClient(
		hls_client.NewRequester(
			testutils.TgAddrs[24],
			&http.Client{Timeout: time.Second},
		),
	)

	if err := client.GetIndex(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := client.GetMessageLimit(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := client.PushMessage(context.Background(), "abc", "hello, world!"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.LoadMessages(context.Background(), "abc", 0, 1, true); err != nil {
		t.Fatal(err)
	}
	if _, err := client.CountMessages(context.Background(), "abc"); err != nil {
		t.Fatal(err)
	}

	msg := dto.NewMessage(true, "hello, world!", time.Now())
	go func() {
		time.Sleep(100 * time.Millisecond)
		msgBroker.Produce("abc", msg)
	}()
	gotMsg, err := client.ListenChat(context.Background(), "abc", "sid")
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
	pMsgBroker message.IMessageBroker,
	pAddress string,
) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc(
		hls_settings.CHandleIndexPath,
		HandleIndexAPI(pLogger),
	) // GET

	mux.HandleFunc(
		hls_settings.CHandleChatMessagePath,
		HandleChatMessageAPI(pCtx, pLogger, pConfig, pHlkClient, pDatabase),
	) // POST

	mux.HandleFunc(
		hls_settings.CHandleChatHistoryLoadPath,
		HandleChatHistoryLoadAPI(pCtx, pLogger, pConfig, pHlkClient, pDatabase),
	) // GET

	mux.HandleFunc(
		hls_settings.CHandleChatHistorySizePath,
		HandleChatHistorySizeAPI(pCtx, pLogger, pConfig, pHlkClient, pDatabase),
	) // GET

	mux.HandleFunc(
		hls_settings.CHandleChatSubscribePath,
		HandleChatSubscribeAPI(pCtx, pLogger, pMsgBroker),
	) // GET

	buildSettings := build.GetSettings()
	return &http.Server{ // nolint: gosec
		Addr:         pAddress,
		Handler:      http.TimeoutHandler(mux, buildSettings.GetHttpCallbackTimeout(), "handle timeout"),
		ReadTimeout:  buildSettings.GetHttpCallbackTimeout(),
		WriteTimeout: buildSettings.GetHttpCallbackTimeout(),
	}
}
