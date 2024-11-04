package handler

import (
	"net/http"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/helpers/loader/pkg/app/config"
	"github.com/number571/hidden-lake/internal/helpers/loader/pkg/settings"
	testutils "github.com/number571/hidden-lake/test/utils"
)

const (
	tcMessageSize = (8 << 10)
	tcWorkSize    = 10
	tcCapacity    = 16
	tcNetworkKey  = "_"
)

var (
	tgProducer = testutils.TgAddrs[25]
	tgConsumer = testutils.TgAddrs[26]
	tgTService = testutils.TgAddrs[27]
)

func testRunService(addr string) *http.Server {
	mux := http.NewServeMux()

	cfg := &config.SConfig{
		FSettings: &config.SConfigSettings{
			FMessagesCapacity: tcCapacity,
			FWorkSizeBits:     tcWorkSize,
			FNetworkKey:       tcNetworkKey,
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
