package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/anonymity/queue"
	"github.com/number571/go-peer/pkg/network/conn"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/service/pkg/app/config"
	pkg_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/closer"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	hiddenlake_network "github.com/number571/hidden-lake/pkg/network"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"
)

const (
	tcMessageSize   = (8 << 10)
	tcWorkSize      = 10
	tcQueuePeriod   = 5_000
	tcFetchTimeout  = 30_000
	tcQueueCapacity = 32
	tcMaxConnects   = 16
	tcCapacity      = 1024
)

var (
	tgPrivKey1 = asymmetric.NewPrivKey()
	tgPrivKey2 = asymmetric.NewPrivKey()
	tgPrivKey3 = asymmetric.NewPrivKey()
)

const (
	tcServiceAddressInHLS = "hidden-echo-service"
	tcPathDBTemplate      = "database_test_%d.db"
	tcPathConfigTemplate  = "config_test_%d.yml"
)

var (
	tcConfig = fmt.Sprintf(`settings:
  message_size_bytes: 8192
  work_size_bits: 22
  fetch_timeout_ms: 60000
  queue_period_ms: 1000
  network_key: test
address:
  tcp: test_address_tcp
  http: test_address_http
connections:
  - test_connect1
  - test_connect2
  - test_connect3
friends:
  test_recvr: %s
  test_name1: %s
services:
  test_service1: test_address1
  test_service2: test_address2
  test_service3: test_address3
`,
		tgPrivKey1.GetPubKey().ToString(),
		tgPrivKey2.GetPubKey().ToString(),
	)
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

func testStartServerHTTP(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/echo", testEchoPage)

	srv := &http.Server{
		Addr:        addr,
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
		ReadTimeout: time.Second,
	}

	go func() { _ = srv.ListenAndServe() }()

	return srv
}

func testEchoPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		FMessage string `json:"message"`
	}

	var resp struct {
		FEcho  string `json:"echo"`
		FError int    `json:"error"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		resp.FError = 1
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	resp.FEcho = req.FMessage
	_ = json.NewEncoder(w).Encode(resp)
}

func testAllCreate(cfgPath, dbPath, srvAddr string) (config.IWrapper, anonymity.INode, context.Context, context.CancelFunc, *http.Server) {
	wcfg := testNewWrapper(cfgPath)
	node, ctx, cancel := testRunNewNode(dbPath, "")
	srvc := testRunService(ctx, wcfg, node, srvAddr)
	time.Sleep(200 * time.Millisecond)
	return wcfg, node, ctx, cancel, srvc
}

func testAllFree(node anonymity.INode, cancel context.CancelFunc, srv *http.Server, pathCfg, pathDB string) {
	defer func() {
		os.RemoveAll(pathDB)
		os.RemoveAll(pathCfg)
	}()
	cancel()
	_ = closer.CloseAll([]io.Closer{
		srv,
		node.GetKVDatabase(),
		node.GetNetworkNode(),
	})
}

func testRunService(ctx context.Context, wcfg config.IWrapper, node anonymity.INode, addr string) *http.Server {
	mux := http.NewServeMux()

	logger := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	cfg := wcfg.GetConfig()

	hlNode := hiddenlake_network.NewRawHiddenLakeNode(
		node,
		func() []string { return nil },
		HandleServiceTCP(cfg, logger),
	)

	mux.HandleFunc(pkg_settings.CHandleIndexPath, HandleIndexAPI(logger))
	mux.HandleFunc(pkg_settings.CHandleConfigSettingsPath, HandleConfigSettingsAPI(wcfg, logger, node))
	mux.HandleFunc(pkg_settings.CHandleConfigConnectsPath, HandleConfigConnectsAPI(ctx, wcfg, logger, node))
	mux.HandleFunc(pkg_settings.CHandleConfigFriendsPath, HandleConfigFriendsAPI(wcfg, logger, node))
	mux.HandleFunc(pkg_settings.CHandleNetworkOnlinePath, HandleNetworkOnlineAPI(logger, node))
	mux.HandleFunc(pkg_settings.CHandleNetworkRequestPath, HandleNetworkRequestAPI(ctx, cfg, logger, hlNode))
	mux.HandleFunc(pkg_settings.CHandleServicePubKeyPath, HandleServicePubKeyAPI(logger, node))

	srv := &http.Server{
		Addr:        addr,
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
		ReadTimeout: time.Second,
	}

	go func() { _ = srv.ListenAndServe() }()
	return srv
}

func testNewWrapper(cfgPath string) config.IWrapper {
	_ = os.WriteFile(cfgPath, []byte(tcConfig), 0o600)
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		panic(err)
	}
	return config.NewWrapper(cfg)
}

func testRunNewNode(dbPath, addr string) (anonymity.INode, context.Context, context.CancelFunc) {
	os.RemoveAll(dbPath)
	node := testNewNode(dbPath, addr).HandleFunc(build.GSettings.FProtoMask.FService, nil)
	ctx, cancel := context.WithCancel(context.Background())
	go func() { _ = node.Run(ctx) }()
	return node, ctx, cancel
}

func testNewNode(dbPath, addr string) anonymity.INode {
	db, err := database.NewKVDatabase(dbPath)
	if err != nil {
		panic(err)
	}
	networkMask := uint32(1)
	node := anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FServiceName:  "TEST",
			FFetchTimeout: time.Minute,
		}),
		logger.NewLogger(
			logger.NewSettings(&logger.SSettings{}),
			func(_ logger.ILogArg) string { return "" },
		),
		db,
		testNewNetworkNode(addr),
		queue.NewQBProblemProcessor(
			queue.NewSettings(&queue.SSettings{
				FMessageConstructSettings: net_message.NewConstructSettings(&net_message.SConstructSettings{
					FSettings: net_message.NewSettings(&net_message.SSettings{
						FWorkSizeBits: tcWorkSize,
					}),
				}),
				FNetworkMask:  networkMask,
				FQueuePeriod:  500 * time.Millisecond,
				FConsumersCap: 1,
				FQueuePoolCap: [2]uint64{tcQueueCapacity, tcQueueCapacity},
			}),
			client.NewClient(
				tgPrivKey1,
				tcMessageSize,
			),
		),
		asymmetric.NewMapPubKeys(),
	)
	return node
}

func testNewNetworkNode(addr string) network.INode {
	return network.NewNode(
		network.NewSettings(&network.SSettings{
			FAddress:      addr,
			FMaxConnects:  tcMaxConnects,
			FReadTimeout:  time.Minute,
			FWriteTimeout: time.Minute,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FMessageSettings: net_message.NewSettings(&net_message.SSettings{
					FWorkSizeBits: tcWorkSize,
				}),
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
	fServiceAddr string
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
func (p *tsConfig) GetConnections() []string { return nil }
func (p *tsConfig) GetService(s string) (string, bool) {
	if s == "hidden-some-host-ok" {
		return p.fServiceAddr, true
	}
	if s == "hidden-some-host-failed" {
		return "localhost:99999", true
	}
	return "", false
}

var (
	_ hiddenlake_network.IHiddenLakeNode = &tsHLNode{}
)

type tsHLNode struct {
	tsNode *tsNode
}

func newTsHiddenLakeNode(tsNode *tsNode) *tsHLNode {
	return &tsHLNode{tsNode: tsNode}
}

func (p *tsHLNode) Run(context.Context) error      { return nil }
func (p *tsHLNode) GetOriginNode() anonymity.INode { return p.tsNode }

func (p *tsHLNode) SendRequest(ctx context.Context, k asymmetric.IPubKey, r request.IRequest) error {
	return p.tsNode.SendPayload(ctx, k, payload.NewPayload64(1, r.ToBytes()))
}

func (p *tsHLNode) FetchRequest(
	ctx context.Context,
	k asymmetric.IPubKey,
	r request.IRequest,
) (response.IResponse, error) {
	b, err := p.tsNode.FetchPayload(ctx, k, payload.NewPayload32(1, r.ToBytes()))
	if err != nil {
		return nil, err
	}
	return response.LoadResponse(b)
}

var (
	_ anonymity.INode = &tsNode{}
)

func newTsNode(pConnectionsOK, pFetchOK, pSendOK, pLoadResponseOK bool) *tsNode {
	return &tsNode{
		fConnectionsOK:  pConnectionsOK,
		fFetchOK:        pFetchOK,
		fSendOK:         pSendOK,
		fLoadResponseOK: pLoadResponseOK,
	}
}

type tsNode struct {
	fConnectionsOK  bool
	fFetchOK        bool
	fSendOK         bool
	fLoadResponseOK bool
}

func (p *tsNode) Run(context.Context) error                              { return nil }
func (p *tsNode) HandleFunc(uint32, anonymity.IHandlerF) anonymity.INode { return p }

func (p *tsNode) GetLogger() logger.ILogger {
	return logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string {
			return ""
		},
	)
}
func (p *tsNode) GetSettings() anonymity.ISettings {
	return anonymity.NewSettings(&anonymity.SSettings{
		FServiceName:  "_",
		FFetchTimeout: time.Second,
	})
}
func (p *tsNode) GetKVDatabase() database.IKVDatabase { return nil }
func (p *tsNode) GetNetworkNode() network.INode       { return &tsNetworkNode{p.fConnectionsOK} }
func (p *tsNode) GetMessageQueue() queue.IQBProblemProcessor {
	return queue.NewQBProblemProcessor(
		queue.NewSettings(&queue.SSettings{
			FMessageConstructSettings: net_message.NewConstructSettings(&net_message.SConstructSettings{
				FSettings: net_message.NewSettings(&net_message.SSettings{}),
			}),
			FQueuePeriod:  5_000,
			FConsumersCap: 1,
			FQueuePoolCap: [2]uint64{16, 16},
		}),
		client.NewClient(asymmetric.NewPrivKey(), 8192),
	)
}

func (p *tsNode) GetMapPubKeys() asymmetric.IMapPubKeys { return asymmetric.NewMapPubKeys() }

func (p *tsNode) SendPayload(context.Context, asymmetric.IPubKey, payload.IPayload64) error {
	if !p.fSendOK {
		return errors.New("some error") // nolint: err113
	}
	return nil
}
func (p *tsNode) FetchPayload(context.Context, asymmetric.IPubKey, payload.IPayload32) ([]byte, error) {
	if !p.fFetchOK {
		return nil, errors.New("some error") // nolint: err113
	}
	if !p.fLoadResponseOK {
		return []byte{123}, nil
	}
	return response.NewResponseBuilder().WithCode(200).Build().ToBytes(), nil
}

var (
	_ network.INode = &tsNetworkNode{}
)

type tsNetworkNode struct {
	fConnectionsOK bool
}

func (p *tsNetworkNode) Close() error                                       { return nil }
func (p *tsNetworkNode) Listen(context.Context) error                       { return nil }
func (p *tsNetworkNode) HandleFunc(uint32, network.IHandlerF) network.INode { return nil }

func (p *tsNetworkNode) GetSettings() network.ISettings {
	return network.NewSettings(&network.SSettings{
		FConnSettings: conn.NewSettings(&conn.SSettings{
			FLimitMessageSizeBytes: 1,
			FWaitReadTimeout:       time.Second,
			FDialTimeout:           time.Second,
			FReadTimeout:           time.Second,
			FWriteTimeout:          time.Second,
			FMessageSettings: net_message.NewSettings(&net_message.SSettings{
				FWorkSizeBits: 1,
				FNetworkKey:   "_",
			}),
		}),
		FMaxConnects:  1,
		FReadTimeout:  time.Second,
		FWriteTimeout: time.Second,
	})
}

func (p *tsNetworkNode) GetCacheSetter() cache.ICacheSetter { return nil }

func (p *tsNetworkNode) GetConnections() map[string]conn.IConn {
	return map[string]conn.IConn{
		"127.0.0.1:9999": nil,
	}
}
func (p *tsNetworkNode) AddConnection(context.Context, string) error { return nil }
func (p *tsNetworkNode) DelConnection(string) error {
	if !p.fConnectionsOK {
		return errors.New("some error") // nolint: err113
	}
	return nil
}

func (p *tsNetworkNode) BroadcastMessage(context.Context, net_message.IMessage) error { return nil }
