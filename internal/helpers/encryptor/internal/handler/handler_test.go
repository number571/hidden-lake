package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/app/config"
	"github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/settings"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

const (
	tcMessageSize = (8 << 10)
	tcNetworkKey  = "_"
	tcHead        = uint32(123)
	tcBody        = "hello, world!"
	tcWorkSize    = 10
)

const (
	tcPathConfigTemplate = "config_test_%d.yml"
)

var (
	tgPrivKey1 = asymmetric.NewPrivKey()
	tgPrivKey2 = asymmetric.NewPrivKey()
	tgPrivKey3 = asymmetric.NewPrivKey()
)

var (
	tcConfig = fmt.Sprintf(`settings:
  message_size_bytes: %d
  work_size_bits: %d
  network_key: %s
address:
  http: test_address_http
friends:
  test_recvr: %s
  test_name1: %s
`,
		tcMessageSize,
		tcWorkSize,
		tcNetworkKey,
		tgPrivKey1.GetPubKey().ToString(),
		tgPrivKey2.GetPubKey().ToString(),
	)
)

func testNewWrapper(cfgPath string) config.IWrapper {
	_ = os.WriteFile(cfgPath, []byte(tcConfig), 0o600)
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		panic(err)
	}
	return config.NewWrapper(cfg)
}

func testInitFriends(pCfg config.IConfig) asymmetric.IMapPubKeys {
	f2f := asymmetric.NewMapPubKeys()
	for _, pubKey := range pCfg.GetFriends() {
		f2f.SetPubKey(pubKey)
	}
	return f2f
}

func testRunService(cfgPath, addr string) (config.IWrapper, *http.Server) {
	mux := http.NewServeMux()

	wcfg := testNewWrapper(cfgPath)
	cfg := wcfg.GetConfig()

	friends := testInitFriends(cfg)

	logger := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	client := client.NewClient(tgPrivKey1, tcMessageSize)

	mux.HandleFunc(settings.CHandleIndexPath, HandleIndexAPI(logger))
	mux.HandleFunc(settings.CHandleMessageEncryptPath, HandleMessageEncryptAPI(cfg, logger, client, 1))
	mux.HandleFunc(settings.CHandleMessageDecryptPath, HandleMessageDecryptAPI(cfg, logger, client, friends))
	mux.HandleFunc(settings.CHandleServicePubKeyPath, HandleServicePubKeyAPI(logger, client.GetPrivKey().GetPubKey()))
	mux.HandleFunc(settings.CHandleConfigSettingsPath, HandleConfigSettingsAPI(cfg, logger))
	mux.HandleFunc(settings.CHandleConfigFriendsPath, HandleConfigFriendsAPI(wcfg, logger, friends))

	srv := &http.Server{
		Addr:        addr,
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
		ReadTimeout: time.Second,
	}

	go func() { _ = srv.ListenAndServe() }()
	return wcfg, srv
}

func testNetworkMessageSettings() net_message.ISettings {
	return net_message.NewSettings(&net_message.SSettings{
		FNetworkKey:   tcNetworkKey,
		FWorkSizeBits: tcWorkSize,
	})
}

type tsWrapper struct {
	fEditorOK bool
}

func newTsWrapper(pEditorOK bool) *tsWrapper {
	return &tsWrapper{
		fEditorOK: pEditorOK,
	}
}

func (p *tsWrapper) GetConfig() config.IConfig { return &tsConfig{} }
func (p *tsWrapper) GetEditor() config.IEditor { return &tsEditor{p.fEditorOK} }

type tsEditor struct {
	fEditorOK bool
}

func (p *tsEditor) UpdateConnections([]string) error {
	if !p.fEditorOK {
		return errors.New("some error") // nolint: err113
	}
	return nil
}
func (p *tsEditor) UpdateFriends(map[string]asymmetric.IPubKey) error {
	if !p.fEditorOK {
		return errors.New("some error") // nolint: err113
	}
	return nil
}

type tsConfig struct {
}

func (p *tsConfig) GetSettings() config.IConfigSettings {
	return &config.SConfigSettings{}
}

func (p *tsConfig) GetLogging() std_logger.ILogging { return nil }
func (p *tsConfig) GetAddress() config.IAddress     { return nil }
func (p *tsConfig) GetFriends() map[string]asymmetric.IPubKey {
	return map[string]asymmetric.IPubKey{
		"abc": tgPrivKey2.GetPubKey(),
	}
}
func (p *tsConfig) GetConnections() []string         { return nil }
func (p *tsConfig) GetService(string) (string, bool) { return "", false }

var (
	_ client.IClient = &tsClient{}
)

func newTsClient(pEncryptOK bool) *tsClient {
	return &tsClient{
		fEncryptOK: pEncryptOK,
	}
}

type tsClient struct {
	fEncryptOK bool
}

func (p *tsClient) GetMessageSize() uint64  { return 8192 }
func (p *tsClient) GetPayloadLimit() uint64 { return 4000 }

func (p *tsClient) GetPrivKey() asymmetric.IPrivKey { return tgPrivKey1 }

func (p *tsClient) EncryptMessage(asymmetric.IPubKey, []byte) ([]byte, error) {
	if !p.fEncryptOK {
		return nil, errors.New("some error") // nolint: err113
	}
	return []byte{1}, nil
}
func (p *tsClient) DecryptMessage(asymmetric.IMapPubKeys, []byte) (asymmetric.IPubKey, []byte, error) {
	return tgPrivKey2.GetPubKey(), []byte{1}, nil
}
