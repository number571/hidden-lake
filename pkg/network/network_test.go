package network

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	gopeer_adapters "github.com/number571/go-peer/pkg/anonymity/adapters"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/pkg/adapters"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SAppError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestPanicNode(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	_ = NewHiddenLakeNode(
		NewSettings(&SSettings{
			FAdapterSettings: adapters.NewSettings(&adapters.SSettings{
				FMessageSizeBytes: 4570,
			}),
			FQBPSettings: &SQBPSettings{
				FQueuePeriod:  time.Second,
				FFetchTimeout: time.Second,
			},
		}),
		asymmetric.NewPrivKey(),
		&tsDatabase{},
		adapters.NewRunnerAdapter(
			gopeer_adapters.NewAdapterByFuncs(
				func(context.Context, layer1.IMessage) error { return nil },
				func(context.Context) (layer1.IMessage, error) { return nil, nil },
			),
			func(context.Context) error { return nil },
		),
		func(_ context.Context, _ asymmetric.IPubKey, _ request.IRequest) (response.IResponse, error) {
			return nil, nil
		},
	)
}

func TestSettings(t *testing.T) {
	t.Parallel()

	defaultNetwork, _ := build.GetNetwork(build.CDefaultNetwork)
	sett := NewSettings(nil)

	if sett.GetAdapterSettings().GetMessageSizeBytes() != defaultNetwork.FMessageSizeBytes {
		t.Error("got invalid message size by default settings")
		return
	}

	if sett.GetAdapterSettings().GetWorkSizeBits() != defaultNetwork.FWorkSizeBits {
		t.Error("got invalid message size by default settings")
		return
	}

	if sett.GetFetchTimeout() != CDefaultFetchTimeout {
		t.Error("got invalid fetch timeout by default settings")
		return
	}

	if sett.GetQueuePeriod() != CDefaultQueuePeriod {
		t.Error("got invalid queue period by default settings")
		return
	}

	if sett.GetQBPConsumers() != CDefaultQBPConsumers {
		t.Error("got invalid qbp_consumers by default")
		return
	}

	if sett.GetPowParallel() != CDefaultPowParallel {
		t.Error("got invalid pow_parallel by default")
		return
	}

	if sett.GetServiceName() != "_" {
		t.Error("got invalid service_name by default")
		return
	}

	sett.GetLogger().PushInfo("___")
}

type tsDatabase struct{}

func (p *tsDatabase) Close() error               { return nil }
func (p *tsDatabase) Set([]byte, []byte) error   { return nil }
func (p *tsDatabase) Get([]byte) ([]byte, error) { return nil, nil }
func (p *tsDatabase) Del([]byte) error           { return nil }

func TestHiddenLakeNode(t *testing.T) {
	t.Parallel()

	msgChan1 := make(chan layer1.IMessage)
	msgChan2 := make(chan layer1.IMessage)

	node1 := testNewHiddenLakeNode("node1.db", msgChan2, msgChan1)
	node1PubKey := node1.GetOriginNode().GetQBProcessor().GetClient().GetPrivKey().GetPubKey()
	defer func() { _ = os.Remove("node1.db") }()

	node2 := testNewHiddenLakeNode("node2.db", msgChan1, msgChan2)
	node2PubKey := node2.GetOriginNode().GetQBProcessor().GetClient().GetPrivKey().GetPubKey()
	defer func() { _ = os.Remove("node2.db") }()

	node1.GetOriginNode().GetMapPubKeys().SetPubKey(node2PubKey)
	node2.GetOriginNode().GetMapPubKeys().SetPubKey(node1PubKey)

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
		t.Error(err)
		return
	}

	rsp, err := node1.FetchRequest(
		ctx,
		node2PubKey,
		request.NewRequestBuilder().WithMethod(http.MethodPut).Build(),
	)
	if err != nil {
		t.Error(err)
		return
	}
	if rsp.GetCode() != http.StatusAccepted {
		t.Error("got invalid status code")
		return
	}
}

func testNewHiddenLakeNode(dbPath string, outMsgChan, inMsgChan chan layer1.IMessage) IHiddenLakeNode {
	return NewHiddenLakeNode(
		NewSettings(&SSettings{
			FAdapterSettings: adapters.NewSettings(&adapters.SSettings{
				FMessageSizeBytes: 8 << 10,
			}),
			FQBPSettings: &SQBPSettings{
				FQueuePeriod:  time.Second,
				FFetchTimeout: time.Minute,
			},
		}),
		asymmetric.NewPrivKey(),
		func() database.IKVDatabase {
			db, err := database.NewKVDatabase(dbPath)
			if err != nil {
				panic(err)
			}
			return db
		}(),
		adapters.NewRunnerAdapter(
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
			_ asymmetric.IPubKey,
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
}
