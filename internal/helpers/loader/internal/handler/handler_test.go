package handler

import (
	"net/http"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	testutils "github.com/number571/go-peer/test/utils"
	"github.com/number571/hidden-lake/internal/helpers/loader/internal/config"
	"github.com/number571/hidden-lake/internal/helpers/loader/pkg/settings"
)

var (
	tgProducer = testutils.TgAddrs[42]
	tgConsumer = testutils.TgAddrs[43]
	tgTService = testutils.TgAddrs[44]
)

func testRunService(addr string) *http.Server {
	mux := http.NewServeMux()

	cfg := &config.SConfig{
		FSettings: &config.SConfigSettings{
			FMessagesCapacity: testutils.TCCapacity,
			FWorkSizeBits:     testutils.TCWorkSize,
			FNetworkKey:       testutils.TCNetworkKey,
		},
		FProducers: []string{tgProducer},
		FConsumers: []string{tgConsumer},
	}

	logger := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	mux.HandleFunc(settings.CHandleIndexPath, HandleIndexAPI(logger))
	mux.HandleFunc(settings.CHandleNetworkTransferPath, HandleNetworkTransferAPI(cfg, logger))
	mux.HandleFunc(settings.CHandleConfigSettings, HandleConfigSettingsAPI(cfg, logger))

	srv := &http.Server{
		Addr:        addr,
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
		ReadTimeout: time.Second,
	}

	go func() { _ = srv.ListenAndServe() }()
	return srv
}
