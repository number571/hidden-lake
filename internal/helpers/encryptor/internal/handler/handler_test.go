package handler

import (
	"net/http"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	net_message "github.com/number571/go-peer/pkg/network/message"
	testutils "github.com/number571/go-peer/test/utils"
	"github.com/number571/hidden-lake/internal/helpers/encryptor/internal/config"
	"github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/settings"
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
			FMessageSizeBytes: testutils.TCMessageSize,
			FWorkSizeBits:     testutils.TCWorkSize,
			FNetworkKey:       testutils.TCNetworkKey,
		},
	}

	logger := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	client := client.NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: testutils.TCMessageSize,
			FEncKeySizeBytes:  asymmetric.CKEncSize,
		}),
		tgPrivKey,
	)

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
		FNetworkKey:   testutils.TCNetworkKey,
		FWorkSizeBits: testutils.TCWorkSize,
	})
}
