package network

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	gopeer_anonymity "github.com/number571/go-peer/pkg/anonymity/qb"
	gopeer_adapters "github.com/number571/go-peer/pkg/anonymity/qb/adapters"
	"github.com/number571/go-peer/pkg/anonymity/qb/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2/hybrid"
	ssymm "github.com/number571/go-peer/pkg/crypto/scheme/layer2/symmetric"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/connkeeper"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/pkg/network/adapters"
	"github.com/number571/hidden-lake/pkg/network/adapters/tcp"
	"github.com/number571/hidden-lake/pkg/network/request"
	"github.com/number571/hidden-lake/pkg/network/response"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestSettings(t *testing.T) {
	t.Parallel()

	defaultNetwork, _ := build.GetNetwork(build.CDefaultNetwork)
	sett := NewSettings(nil)

	if sett.GetAdapterSettings().GetMessageSizeBytes() != defaultNetwork.FMessageSizeBytes {
		t.Fatal("got invalid message size by default settings")
	}

	if sett.GetAdapterSettings().GetWorkSizeBits() != 0 {
		t.Fatal("got invalid work size bits by default value")
	}

	if sett.GetFetchTimeout() != CDefaultFetchTimeout {
		t.Fatal("got invalid fetch timeout by default settings")
	}

	if sett.GetQBPSettings().GetGeneratePeriod() != CDefaultQueuePeriod {
		t.Fatal("got invalid queue period by default settings")
	}

	if sett.GetQBPSettings().GetNumberOfConsumers() != 1 {
		t.Fatal("got invalid qbp_consumers by default")
	}

	if sett.GetPowParallel() != 1 {
		t.Fatal("got invalid pow_parallel by default")
	}

	if sett.GetFmtAppName() != "_" {
		t.Fatal("got invalid service_name by default")
	}

	sett.GetLogger().PushInfo("___")
}

func TestHiddenLakeNodeHybridScheme(t *testing.T) {
	t.Parallel()

	msgChan1 := make(chan layer1.IMessage)
	msgChan2 := make(chan layer1.IMessage)

	if err := testNewHiddenLakeNodeHybridSchemeInvalidPayloadSize("node1_hybrid_x.db", msgChan2, msgChan1); err == nil {
		t.Fatal("success create node with invalid payload size")
	}
	defer func() { _ = os.Remove("node1_hybrid_x.db") }()

	node1, node1PubKey := testNewHiddenLakeNodeHybridScheme("node1_hybrid.db", msgChan2, msgChan1)
	defer func() { _ = os.Remove("node1_hybrid.db") }()

	node2, node2PubKey := testNewHiddenLakeNodeHybridScheme("node2_hybrid.db", msgChan1, msgChan2)
	defer func() { _ = os.Remove("node2_hybrid.db") }()

	node1.GetOriginNode().GetKeysContainer().Add(node2PubKey)
	node2.GetOriginNode().GetKeysContainer().Add(node1PubKey)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() { _ = node1.Run(ctx) }()
	go func() { _ = node2.Run(ctx) }()

	err := node1.SendRequest(
		ctx,
		node2PubKey,
		request.NewRequestBuilder().WithMethod(http.MethodPost).Build(),
	)
	if err != nil {
		t.Fatal(err)
	}

	rsp, err := node1.FetchRequest(
		ctx,
		node2PubKey,
		request.NewRequestBuilder().WithMethod(http.MethodPut).Build(),
	)
	if err != nil {
		t.Fatal(err)
	}
	if rsp.GetCode() != http.StatusAccepted {
		t.Fatal("got invalid status code")
	}

	err = node1.SendRequest(
		ctx,
		node2PubKey,
		request.NewRequestBuilder().WithMethod(http.MethodPost).WithBody(random.NewRandom().GetBytes(10<<10)).Build(),
	)
	if err == nil {
		t.Fatal("success send invalid request")
	}

	_, err = node1.FetchRequest(
		ctx,
		node2PubKey,
		request.NewRequestBuilder().WithMethod(http.MethodPost).WithBody(random.NewRandom().GetBytes(10<<10)).Build(),
	)
	if err == nil {
		t.Fatal("success fetch invalid request")
	}
}

func TestHiddenLakeNodeSymmetricScheme(t *testing.T) {
	t.Parallel()

	msgChan1 := make(chan layer1.IMessage)
	msgChan2 := make(chan layer1.IMessage)

	if err := testNewHiddenLakeNodeSymmetricSchemeInvalidPayloadSize("node1_symmetric_x.db", msgChan2, msgChan1); err == nil {
		t.Fatal("success create node with invalid payload size")
	}
	defer func() { _ = os.Remove("node1_symmetric_x.db") }()

	node1, pKey1 := testNewHiddenLakeNodeSymmetricScheme("node1_symmetric.db", msgChan2, msgChan1)
	defer func() { _ = os.Remove("node1_symmetric.db") }()

	node2, pKey2 := testNewHiddenLakeNodeSymmetricScheme("node2_symmetric.db", msgChan1, msgChan2)
	defer func() { _ = os.Remove("node2_symmetric.db") }()

	node1.GetOriginNode().GetKeysContainer().Add(pKey1)
	node2.GetOriginNode().GetKeysContainer().Add(pKey2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() { _ = node1.Run(ctx) }()
	go func() { _ = node2.Run(ctx) }()

	err := node1.SendRequest(
		ctx,
		pKey2,
		request.NewRequestBuilder().WithMethod(http.MethodPost).Build(),
	)
	if err != nil {
		t.Fatal(err)
	}

	rsp, err := node1.FetchRequest(
		ctx,
		pKey2,
		request.NewRequestBuilder().WithMethod(http.MethodPut).Build(),
	)
	if err != nil {
		t.Fatal(err)
	}
	if rsp.GetCode() != http.StatusAccepted {
		t.Fatal("got invalid status code")
	}

	err = node1.SendRequest(
		ctx,
		pKey2,
		request.NewRequestBuilder().WithMethod(http.MethodPost).WithBody(random.NewRandom().GetBytes(10<<10)).Build(),
	)
	if err == nil {
		t.Fatal("success send invalid request")
	}

	_, err = node1.FetchRequest(
		ctx,
		pKey2,
		request.NewRequestBuilder().WithMethod(http.MethodPost).WithBody(random.NewRandom().GetBytes(10<<10)).Build(),
	)
	if err == nil {
		t.Fatal("success fetch invalid request")
	}
}

func TestHiddenLakeNodeCustom(t *testing.T) {
	t.Parallel()

	msgChan1 := make(chan layer1.IMessage)
	msgChan2 := make(chan layer1.IMessage)

	node1, _ := testNewHiddenLakeNodeCustom("node1_custom.db", msgChan2, msgChan1)
	defer func() { _ = os.Remove("node1_custom.db") }()

	_node := node1.(*sHiddenLakeNode)
	_node.fOriginNode = &tsNode{fLoadResponseOK: false}

	scheme := node1.GetOriginNode().GetQBProcessor().GetScheme()
	_, err := node1.FetchRequest(
		context.Background(),
		scheme.GetRandomKey(),
		request.NewRequestBuilder().WithMethod(http.MethodPost).WithBody(random.NewRandom().GetBytes(10<<10)).Build(),
	)
	if err == nil {
		t.Fatal("success fetch with invalid response")
	}
	if err := node1.Run(context.Background()); err == nil {
		t.Fatal("success run node with got error from runner")
	}
}

func TestGetQueueProcessor(t *testing.T) {
	t.Parallel()

	settings := NewSettings(&SSettings{
		FAdapterSettings: adapters.NewSettings(&adapters.SSettings{}),
		FQBPSettings: &SQBPSettings{
			FIsDisabled:     true,
			FGeneratePeriod: time.Second,
		},
		FServeSettings: &SServeSettings{
			FFetchTimeout: time.Minute,
		},
	})

	scheme, _ := ssymm.NewScheme(settings.GetAdapterSettings().GetMessageSizeBytes())
	_ = getQueueProcessor(settings, scheme)
}

func testNewHiddenLakeNodeSymmetricSchemeInvalidPayloadSize(dbPath string, outMsgChan, inMsgChan chan layer1.IMessage) error {
	adapterSettings := adapters.NewSettings(&adapters.SSettings{FMessageSizeBytes: ssymm.CMessageHeadSize + 3})
	scheme, _ := ssymm.NewScheme(adapterSettings.GetMessageSizeBytes())
	_, err := testNewHiddenLakeNode(dbPath, outMsgChan, inMsgChan, scheme, adapterSettings)
	return err
}

func testNewHiddenLakeNodeSymmetricScheme(dbPath string, outMsgChan, inMsgChan chan layer1.IMessage) (IHiddenLakeNode, layer2.IParticipantKey) {
	key := []byte("secure_shared_secret_for_friends")
	adapterSettings := adapters.NewSettings(&adapters.SSettings{FMessageSizeBytes: 128})
	scheme, _ := ssymm.NewScheme(adapterSettings.GetMessageSizeBytes())
	node, _ := testNewHiddenLakeNode(dbPath, outMsgChan, inMsgChan, scheme, adapterSettings)
	return node, symmetric.NewCipherGCM(key)
}

func testNewHiddenLakeNodeCustom(dbPath string, outMsgChan, inMsgChan chan layer1.IMessage) (IHiddenLakeNode, layer2.IParticipantKey) {
	key := []byte("secure_shared_secret_for_friends")
	adapterSettings := adapters.NewSettings(&adapters.SSettings{FMessageSizeBytes: 128})
	scheme, _ := ssymm.NewScheme(adapterSettings.GetMessageSizeBytes())
	node, _ := testNewHiddenLakeNode(dbPath, outMsgChan, inMsgChan, scheme, adapterSettings)
	return node, symmetric.NewCipherGCM(key)
}

func testNewHiddenLakeNodeHybridSchemeInvalidPayloadSize(dbPath string, outMsgChan, inMsgChan chan layer1.IMessage) error {
	adapterSettings := adapters.NewSettings(&adapters.SSettings{FMessageSizeBytes: hybrid.CMessageHeadSize + 3})
	privKey := asymmetric.NewPrivKey()
	scheme, _ := hybrid.NewScheme(privKey, adapterSettings.GetMessageSizeBytes())
	_, err := testNewHiddenLakeNode(dbPath, outMsgChan, inMsgChan, scheme, adapterSettings)
	return err
}

func testNewHiddenLakeNodeHybridScheme(dbPath string, outMsgChan, inMsgChan chan layer1.IMessage) (IHiddenLakeNode, layer2.IParticipantKey) {
	adapterSettings := adapters.NewSettings(&adapters.SSettings{FMessageSizeBytes: 8 << 10})
	privKey := asymmetric.NewPrivKey()
	scheme, _ := hybrid.NewScheme(privKey, adapterSettings.GetMessageSizeBytes())
	node, _ := testNewHiddenLakeNode(dbPath, outMsgChan, inMsgChan, scheme, adapterSettings)
	return node, privKey.GetPubKey()
}

func testNewHiddenLakeNode(dbPath string, outMsgChan, inMsgChan chan layer1.IMessage, scheme layer2.IScheme, adapterSettings adapters.ISettings) (IHiddenLakeNode, error) {
	node, err := NewHiddenLakeNode(
		NewSettings(&SSettings{
			FAdapterSettings: adapterSettings,
			FQBPSettings: &SQBPSettings{
				FGeneratePeriod: time.Second,
			},
			FServeSettings: &SServeSettings{
				FFetchTimeout: time.Minute,
			},
		}),
		scheme,
		layer2.NewKeysContainer(),
		func() database.IKVDatabase {
			db, err := database.NewKVDatabase(dbPath)
			if err != nil {
				panic(err)
			}
			return db
		}(),
		newRunnerAdapter(
			gopeer_adapters.NewAdapterByFuncs(
				func(_ context.Context, msg layer1.IMessage) error {
					outMsgChan <- msg
					return nil
				},
				func(_ context.Context) (layer1.IMessage, error) {
					return <-inMsgChan, nil
				},
			),
			func(ctx context.Context) error {
				<-ctx.Done()
				return nil
			},
		),
		func(
			_ context.Context,
			_ layer2.IParticipantKey,
			req request.IRequest,
		) (response.IResponse, error) {
			if req.GetMethod() == http.MethodPost {
				return nil, nil
			}
			if req.GetMethod() == http.MethodPut {
				return response.NewResponseBuilder().WithCode(http.StatusAccepted).Build(), nil
			}
			panic("unknown method")
		},
	)
	if err != nil {
		return nil, err
	}
	return node, nil
}

type sRunnerAdapter struct {
	gopeer_adapters.IAdapter
	fRun func(context.Context) error
}

func newRunnerAdapter(pAdapter gopeer_adapters.IAdapter, pRun func(context.Context) error) adapters.IRunnerAdapter {
	return &sRunnerAdapter{
		IAdapter: pAdapter,
		fRun:     pRun,
	}
}

func (p *sRunnerAdapter) Run(pCtx context.Context) error {
	return p.fRun(pCtx)
}

type tsNode struct {
	fLoadResponseOK bool
}

func (p *tsNode) Run(context.Context) error                                            { return errors.New("some error") } //nolint:err113
func (p *tsNode) HandleFunc(uint32, gopeer_anonymity.IHandlerF) gopeer_anonymity.INode { return p }
func (p *tsNode) GetAdapter() gopeer_adapters.IAdapter {
	return &tsTCPAdapter{}
}

func (p *tsNode) GetLogger() logger.ILogger {
	return logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string {
			return ""
		},
	)
}
func (p *tsNode) GetSettings() gopeer_anonymity.ISettings {
	return gopeer_anonymity.NewSettings(&gopeer_anonymity.SSettings{
		FServiceName:  "_",
		FFetchTimeout: time.Second,
	})
}
func (p *tsNode) GetKVDatabase() database.IKVDatabase { return nil }
func (p *tsNode) GetNetworkNode() network.INode       { return &tsNetworkNode{} }
func (p *tsNode) GetQBProcessor() queue.IQBProblemProcessor {
	return queue.NewQBProblemProcessor(
		queue.NewSettings(&queue.SSettings{
			FMessageConstructSettings: layer1.NewConstructSettings(&layer1.SConstructSettings{
				FSettings: layer1.NewSettings(&layer1.SSettings{}),
			}),
			FQueuePeriod:  5_000,
			FConsumersCap: 1,
			FQueuePoolCap: [2]uint64{16, 16},
		}),
		func() layer2.IScheme {
			scheme, _ := hybrid.NewScheme(asymmetric.NewPrivKey(), 8192)
			return scheme
		}(),
	)
}

func (p *tsNode) GetKeysContainer() layer2.IKeysContainer { return layer2.NewKeysContainer() }

func (p *tsNode) SendPayload(context.Context, layer2.IParticipantKey, []byte) error {
	return nil
}
func (p *tsNode) FetchPayload(context.Context, layer2.IParticipantKey, []byte) ([]byte, error) {
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
func (p *tsNetworkNode) Run(context.Context) error                          { return nil }
func (p *tsNetworkNode) HandleFunc(uint32, network.IHandlerF) network.INode { return nil }

func (p *tsNetworkNode) GetSettings() network.ISettings {
	return network.NewSettings(&network.SSettings{
		FConnSettings: conn.NewSettings(&conn.SSettings{
			FLimitMessageSizeBytes: 1,
			FWaitReadTimeout:       time.Second,
			FDialTimeout:           time.Second,
			FReadTimeout:           time.Second,
			FWriteTimeout:          time.Second,
			FMessageSettings: layer1.NewSettings(&layer1.SSettings{
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

func (p *tsNetworkNode) BroadcastMessage(context.Context, layer1.IMessage) error { return nil }

var (
	_ tcp.ITCPAdapter        = &tsTCPAdapter{}
	_ connkeeper.IConnKeeper = &tsConnKeeper{}
)

type tsTCPAdapter struct {
	fConnectionsOK bool
}
type tsConnKeeper struct {
	fConnectionsOK bool
}

func (p *tsTCPAdapter) Produce(context.Context, layer1.IMessage) error   { return nil }
func (p *tsTCPAdapter) Consume(context.Context) (layer1.IMessage, error) { return nil, nil }
func (p *tsTCPAdapter) WithLogger(_ string, _ logger.ILogger) tcp.ITCPAdapter {
	return p
}
func (p *tsTCPAdapter) GetConnKeeper() connkeeper.IConnKeeper {
	return &tsConnKeeper{p.fConnectionsOK}
}
func (p *tsTCPAdapter) Run(context.Context) error { return nil }

func (p *tsConnKeeper) GetSettings() connkeeper.ISettings { return nil }
func (p *tsConnKeeper) GetNetworkNode() network.INode {
	return &tsNetworkNode{fConnectionsOK: p.fConnectionsOK}
}
func (p *tsConnKeeper) Run(context.Context) error { return nil }
