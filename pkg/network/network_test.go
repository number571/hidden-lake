package network

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	gopeer_adapters "github.com/number571/go-peer/pkg/anonymity/qb/adapters"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2/hybrid"
	ssymm "github.com/number571/go-peer/pkg/crypto/scheme/layer2/symmetric"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/pkg/network/adapters"
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

	if sett.GetAdapterSettings().GetWorkSizeBits() != defaultNetwork.FWorkSizeBits {
		t.Fatal("got invalid message size by default settings")
	}

	if sett.GetFetchTimeout() != CDefaultFetchTimeout {
		t.Fatal("got invalid fetch timeout by default settings")
	}

	if sett.GetQueuePeriod() != CDefaultQueuePeriod {
		t.Fatal("got invalid queue period by default settings")
	}

	if sett.GetQBPConsumers() != 1 {
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
				FQueuePeriod:  time.Second,
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
