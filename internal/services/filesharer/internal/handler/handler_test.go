package handler

import (
	"bytes"
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/app/config"
	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	hls_client "github.com/number571/hidden-lake/pkg/api/services/filesharer/client"
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

	hlkClient := &tsHLKClientWrapper{fClient: newTsHLKClient(0, true)}
	server := testInitInternalServiceHTTP(
		ctx,
		&tsConfig{},
		httpLogger,
		hlkClient,
		"./testdata",
		testutils.TgAddrs[27],
	)
	defer func() { _ = server.Close() }()

	go func() { _ = server.ListenAndServe() }()
	time.Sleep(100 * time.Millisecond)

	client := hls_client.NewClient(
		hls_client.NewRequester(
			testutils.TgAddrs[27],
			&http.Client{Timeout: time.Second},
		),
	)

	if _, err := client.GetIndex(context.Background()); err != nil {
		t.Fatal(err)
	}

	if _, err := client.GetLocalList(context.Background(), "", 0); err != nil {
		t.Fatal(err)
	}
	if _, err := client.GetLocalFileInfo(context.Background(), "", "example.txt"); err != nil {
		t.Fatal(err)
	}
	w1 := &strings.Builder{}
	if err := client.GetLocalFile(w1, context.Background(), "", "example.txt"); err != nil {
		t.Fatal(err)
	}
	if err := client.PutLocalFile(context.Background(), "", "test111.txt", bytes.NewBuffer([]byte("hello"))); err != nil {
		t.Fatal(err)
	}
	if err := client.DelLocalFile(context.Background(), "", "test111.txt"); err != nil {
		t.Fatal(err)
	}

	if _, err := client.GetRemoteFileInfo(context.Background(), "abc", "example.txt", false); err != nil {
		t.Fatal(err)
	}

	hlkClient.fClient = newTsHLKClient(2, true)
	if _, err := client.GetRemoteList(context.Background(), "abc", 0, false); err != nil {
		t.Fatal(err)
	}

	hlkClient.fClient = newTsHLKClient(3, true)
	w2 := &strings.Builder{}
	if _, err := client.GetRemoteFile(w2, context.Background(), "abc", "example.txt", false); err != nil {
		t.Fatal(err)
	}
}

func testInitInternalServiceHTTP(
	pCtx context.Context,
	pConfig config.IConfig,
	pLogger logger.ILogger,
	pHlkClient hlk_client.IClient,
	pPathTo string,
	pAddress string,
) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc(
		hls_settings.CHandleIndexPath,
		HandleIndexAPI(pLogger),
	) // GET

	mux.HandleFunc(
		hls_settings.CHandleRemoteFileInfoPath,
		HandleRemoteFileInfoAPI(pCtx, pLogger, pHlkClient),
	) // GET

	mux.HandleFunc(
		hls_settings.CHandleRemoteFilePath,
		HandleRemoteFileAPI(pCtx, pConfig, pLogger, pHlkClient, pPathTo),
	) // GET

	mux.HandleFunc(
		hls_settings.CHandleRemoteListPath,
		HandleRemoteListAPI(pCtx, pLogger, pHlkClient),
	) // GET

	mux.HandleFunc(
		hls_settings.CHandleLocalListPath,
		HandleLocalListAPI(pCtx, pConfig, pLogger, pHlkClient, pPathTo),
	) // GET

	mux.HandleFunc(
		hls_settings.CHandleLocalFilePath,
		HandleLocalFileAPI(pCtx, pLogger, pHlkClient, pPathTo),
	) // GET, POST, DELETE

	mux.HandleFunc(
		hls_settings.CHandleLocalFileInfoPath,
		HandleLocalFileInfoAPI(pCtx, pLogger, pHlkClient, pPathTo),
	) // GET

	return &http.Server{ // nolint: gosec
		Addr:    pAddress,
		Handler: mux,
	}
}
