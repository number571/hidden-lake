package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/connkeeper"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/internal/helpers/traffic/internal/config"
	"github.com/number571/hidden-lake/internal/helpers/traffic/internal/storage"
	hlt_client "github.com/number571/hidden-lake/internal/helpers/traffic/pkg/client"
	pkg_settings "github.com/number571/hidden-lake/internal/helpers/traffic/pkg/settings"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/closer"
)

const (
	tcMaxConnects    = 16
	tcCapacity       = 32
	tcWorkSize       = 10
	tcNetworkKey     = "_"
	tcMessageSize    = (8 << 10)
	tcHead           = uint32(123)
	tcBody           = "hello, world!"
	databaseTemplate = "database_test_%s.db"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SHandlerError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func testNetworkMessageSettings() net_message.IConstructSettings {
	return net_message.NewConstructSettings(&net_message.SConstructSettings{
		FSettings: net_message.NewSettings(&net_message.SSettings{
			FNetworkKey:   tcNetworkKey,
			FWorkSizeBits: tcWorkSize,
		}),
	})
}

func testAllRun(addr string) (*http.Server, context.CancelFunc, storage.IMessageStorage, hlt_client.IClient) {
	db, err := database.NewKVDatabase(fmt.Sprintf(databaseTemplate, addr))
	if err != nil {
		panic(err)
	}

	stg := storage.NewMessageStorage(
		net_message.NewSettings(&net_message.SSettings{
			FNetworkKey:   tcNetworkKey,
			FWorkSizeBits: tcWorkSize,
		}),
		db,
		cache.NewLRUCache(tcCapacity),
	)

	srv, _, cancel := testRunService(stg, addr, "")

	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			"http://"+addr,
			&http.Client{Timeout: time.Minute},
			testNetworkMessageSettings().GetSettings(),
		),
	)

	time.Sleep(200 * time.Millisecond)
	return srv, cancel, stg, hltClient
}

func testAllFree(addr string, srv *http.Server, cancel context.CancelFunc, db storage.IMessageStorage) {
	cancel()
	defer func() {
		os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))
	}()
	_ = closer.CloseAll([]types.ICloser{srv, db.GetKVDatabase()})
}

func testRunService(stg storage.IMessageStorage, addr string, addrNode string) (*http.Server, connkeeper.IConnKeeper, context.CancelFunc) {
	mux := http.NewServeMux()

	connKeeperSettings := &connkeeper.SSettings{
		FDuration: time.Minute,
		FConnections: func() []string {
			return nil
		},
	}

	if addrNode != "" {
		connKeeperSettings.FConnections = func() []string {
			return []string{addrNode}
		}
	}

	connKeeper := connkeeper.NewConnKeeper(
		connkeeper.NewSettings(connKeeperSettings),
		testNewNetworkNode("").HandleFunc(
			1, // default value
			func(_ context.Context, _ network.INode, _ conn.IConn, _ net_message.IMessage) error {
				// pass response actions
				return nil
			},
		),
	)

	ctx, cancel := context.WithCancel(context.Background())
	go func() { _ = connKeeper.Run(ctx) }()

	cfg := &config.SConfig{
		FSettings: &config.SConfigSettings{
			FMessageSizeBytes:     tcMessageSize,
			FWorkSizeBits:         tcWorkSize,
			FRandMessageSizeBytes: hls_settings.CDefaultRandMessageSizeBytes,
			FNetworkKey:           tcNetworkKey,
			FMessagesCapacity:     tcCapacity,
		},
	}

	node := testNewNetworkNode("")
	logger := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	mux.HandleFunc(pkg_settings.CHandleIndexPath, HandleIndexAPI(logger))
	mux.HandleFunc(pkg_settings.CHandleStoragePointerPath, HandlePointerAPI(stg, logger))
	mux.HandleFunc(pkg_settings.CHandleStorageHashesPath, HandleHashesAPI(stg, logger))
	mux.HandleFunc(pkg_settings.CHandleNetworkMessagePath, HandleMessageAPI(ctx, cfg, stg, logger, logger, node))
	mux.HandleFunc(pkg_settings.CHandleConfigSettings, HandleConfigSettingsAPI(cfg, logger))

	srv := &http.Server{
		Addr:        addr,
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
		ReadTimeout: time.Second,
	}

	go func() { _ = srv.ListenAndServe() }()

	return srv, connKeeper, cancel
}

func testNewClient() client.IClient {
	privKey := asymmetric.NewPrivKeyChain(
		asymmetric.NewKEncPrivKey(),
		asymmetric.NewSignPrivKey(),
	)
	return client.NewClient(privKey, tcMessageSize)
}

func testNewNetworkNode(addr string) network.INode {
	return network.NewNode(
		network.NewSettings(&network.SSettings{
			FAddress:      addr,
			FMaxConnects:  tcMaxConnects,
			FReadTimeout:  time.Minute,
			FWriteTimeout: time.Minute,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FMessageSettings:       testNetworkMessageSettings().GetSettings(),
				FLimitMessageSizeBytes: tcMessageSize,
				FWaitReadTimeout:       time.Hour,
				FDialTimeout:           time.Minute,
				FReadTimeout:           time.Minute,
				FWriteTimeout:          time.Minute,
			}),
		}),
		cache.NewLRUCache(tcCapacity),
	)
}
