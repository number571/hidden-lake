package handler

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/remoter/pkg/app/config"
	"github.com/number571/hidden-lake/internal/applications/remoter/pkg/settings"
	"github.com/number571/hidden-lake/internal/service/pkg/app"
	hls_config "github.com/number571/hidden-lake/internal/service/pkg/app/config"
)

var (
	tgPrivKey1 = asymmetric.NewPrivKey()
	tgPrivKey2 = asymmetric.NewPrivKey()
)

const (
	tcPassword             = "test-password"
	tcPathConfigTemplate   = "config_test_%d.yml"
	tcPathDatabaseTemplate = "database_test_%d.yml"
)

func testRunService(addr string) (config.IConfig, *http.Server) {
	mux := http.NewServeMux()

	cfg := &config.SConfig{
		FSettings: &config.SConfigSettings{
			FPassword:      tcPassword,
			FExecTimeoutMS: 5_000,
		},
		FAddress: &config.SAddress{
			FIncoming: addr,
		},
	}

	logger := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	ctx := context.Background()

	mux.HandleFunc(settings.CExecPath, HandleIncomingExecHTTP(ctx, cfg, logger))

	srv := &http.Server{
		Addr:        addr,
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
		ReadTimeout: time.Second,
	}

	go func() { _ = srv.ListenAndServe() }()
	return cfg, srv
}

func testRunNewNodes(ctx context.Context, httpAddrNode1, tcpAddrNode2, addrService string) {
	runner1 := app.NewApp(
		initConfig(
			"./testdata/node1/hls.yml",
			"",
			httpAddrNode1,
			"",
			tcpAddrNode2,
			tgPrivKey2.GetPubKey(),
		),
		tgPrivKey1,
		"./testdata/node1",
		1,
	)

	runner2 := app.NewApp(
		initConfig(
			"./testdata/node2/hls.yml",
			tcpAddrNode2,
			"",
			addrService,
			"",
			tgPrivKey1.GetPubKey(),
		),
		tgPrivKey2,
		"./testdata/node2",
		1,
	)

	go func() { _ = runner1.Run(ctx) }()
	go func() { _ = runner2.Run(ctx) }()
}

func initConfig(
	cfgPath, tcpAddrNode, httpAddrNode, addrService, conn string,
	pubKey asymmetric.IPubKey,
) hls_config.IConfig {
	os.Remove(cfgPath)

	services := map[string]string{}
	if addrService != "" {
		services[settings.CServiceFullName] = addrService
	}

	connections := []string{}
	if conn != "" {
		connections = []string{conn}
	}

	cfg, err := hls_config.BuildConfig(cfgPath, &hls_config.SConfig{
		FSettings: &hls_config.SConfigSettings{
			FMessageSizeBytes: (8 << 10),
			FWorkSizeBits:     1,
			FFetchTimeoutMS:   10_000,
			FQueuePeriodMS:    500,
		},
		FLogging: []string{},
		FAddress: &hls_config.SAddress{
			FTCP:  tcpAddrNode,
			FHTTP: httpAddrNode,
		},
		FServices:    services,
		FConnections: connections,
		FFriends: map[string]string{
			"test_recv": pubKey.ToString(),
		},
	})
	if err != nil {
		panic(err)
	}
	return cfg
}
