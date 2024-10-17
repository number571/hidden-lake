package handler

import (
	"net/http"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/hidden-lake/internal/helpers/encryptor/internal/config"
	"github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/settings"
)

const (
	tcMessageSize = (8 << 10)
	tcNetworkKey  = "_"
	tcHead        = uint32(123)
	tcBody        = "hello, world!"
	tcWorkSize    = 10
)

var (
	tgPrivKey = asymmetric.NewPrivKeyChain(
		asymmetric.NewKEncPrivKey(),
		asymmetric.NewSignPrivKey(),
	)
)

func testRunService(addr string) *http.Server {
	mux := http.NewServeMux()

	cfg := &config.SConfig{
		FSettings: &config.SConfigSettings{
			FMessageSizeBytes: tcMessageSize,
			FWorkSizeBits:     tcWorkSize,
			FNetworkKey:       tcNetworkKey,
		},
	}

	logger := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	client := client.NewClient(tgPrivKey, tcMessageSize)

	mux.HandleFunc(settings.CHandleIndexPath, HandleIndexAPI(logger))
	mux.HandleFunc(settings.CHandleMessageEncryptPath, HandleMessageEncryptAPI(cfg, logger, client, 1))
	mux.HandleFunc(settings.CHandleMessageDecryptPath, HandleMessageDecryptAPI(cfg, logger, client))
	mux.HandleFunc(settings.CHandleServicePubKeyPath, HandleServicePubKeyAPI(logger, client.GetPrivKeyChain().GetPubKeyChain()))
	mux.HandleFunc(settings.CHandleConfigSettings, HandleConfigSettingsAPI(cfg, logger))

	srv := &http.Server{
		Addr:        addr,
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
		ReadTimeout: time.Second,
	}

	go func() { _ = srv.ListenAndServe() }()
	return srv
}

func testNetworkMessageSettings() net_message.ISettings {
	return net_message.NewSettings(&net_message.SSettings{
		FNetworkKey:   tcNetworkKey,
		FWorkSizeBits: tcWorkSize,
	})
}
