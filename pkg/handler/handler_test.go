package handler

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/anonymity"
	"github.com/number571/go-peer/pkg/anonymity/adapters"
	"github.com/number571/go-peer/pkg/anonymity/queue"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"
)

var (
	tgPubKey = asymmetric.NewPrivKey().GetPubKey()
)

func TestRequestHandler(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	node := &tsNode{}

	if _, err := RequestHandler(testHandler(2))(ctx, node, tgPubKey, []byte{}); err == nil {
		t.Error("success handle with invalid bytes")
		return
	}

	msg := []byte("hello")
	req := request.NewRequestBuilder().WithBody(msg).Build().ToBytes()
	rspb, err := RequestHandler(testHandler(2))(ctx, node, tgPubKey, req)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := response.LoadResponse(rspb)
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(rsp.GetBody(), msg) {
		t.Error("invalid response bytes")
		return
	}

	rspb2, err := RequestHandler(testHandler(1))(ctx, node, tgPubKey, req)
	if err != nil {
		t.Error(err)
		return
	}
	if rspb2 != nil {
		t.Error("get response bytes with return nil")
		return
	}

	if _, err := RequestHandler(testHandler(0))(ctx, node, tgPubKey, req); err == nil {
		t.Error("success handle with invalid response")
		return
	}
}

func testHandler(pMode int) IHandlerF {
	return func(
		_ context.Context,
		pPubKey asymmetric.IPubKey,
		pRequest request.IRequest,
	) (response.IResponse, error) {
		if pMode == 0 {
			return nil, errors.New("some error") // nolint: err113
		}
		if pPubKey.ToString() != tgPubKey.ToString() {
			return nil, errors.New("invalid pub key") // nolint: err113
		}
		if pMode == 1 {
			return nil, nil
		}
		if pMode == 2 {
			return response.NewResponseBuilder().WithBody(pRequest.GetBody()).Build(), nil
		}
		panic("unknown mode")
	}
}

type tsNode struct{}

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
func (p *tsNode) GetNetworkNode() network.INode       { return &tsNetworkNode{} }
func (p *tsNode) GetQBProcessor() queue.IQBProblemProcessor {
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

func (p *tsNode) GetAdapter() adapters.IAdapter         { return nil }
func (p *tsNode) GetMapPubKeys() asymmetric.IMapPubKeys { return asymmetric.NewMapPubKeys() }

func (p *tsNode) SendPayload(context.Context, asymmetric.IPubKey, payload.IPayload64) error {
	return nil
}
func (p *tsNode) FetchPayload(context.Context, asymmetric.IPubKey, payload.IPayload32) ([]byte, error) {
	return response.NewResponseBuilder().WithCode(200).Build().ToBytes(), nil
}

var (
	_ network.INode = &tsNetworkNode{}
)

type tsNetworkNode struct {
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
	return nil
}

func (p *tsNetworkNode) BroadcastMessage(context.Context, net_message.IMessage) error { return nil }
