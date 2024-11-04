package adapters

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	hlt_config "github.com/number571/hidden-lake/internal/helpers/traffic/pkg/config"
	"github.com/number571/hidden-lake/internal/utils/api"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestConsumer(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logging, err := std_logger.LoadLogging([]string{})
	if err != nil {
		t.Error(err)
		return
	}

	httpLogger := std_logger.NewStdLogger(
		logging,
		func(_ logger.ILogArg) string {
			return ""
		},
	)

	hltClient := &tsHLTClient{}
	consumer := &tsConsumer{
		fErr: make(chan error),
		fMsg: make(chan net_message.IMessage),
	}
	go func() {
		err := ConsumeProcessor(ctx, consumer, httpLogger, hltClient, time.Millisecond)
		if err != nil {
			return
		}
	}()

	sett := net_message.NewSettings(&net_message.SSettings{
		FWorkSizeBits: 1,
		FNetworkKey:   "_",
	})

	msg := net_message.NewMessage(
		net_message.NewConstructSettings(&net_message.SConstructSettings{
			FSettings: sett,
			FParallel: 1,
		}),
		payload.NewPayload32(1, []byte("hello, world!")),
	)

	consumer.fMsg <- nil
	consumer.fMsg <- msg
	consumer.fErr <- errors.New("some error") // nolint: err113

	hltClient.fNeedFailed = true
	consumer.fMsg <- msg
	consumer.fMsg <- nil
}

func TestProducer(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logging, err := std_logger.LoadLogging([]string{})
	if err != nil {
		t.Error(err)
		return
	}

	httpLogger := std_logger.NewStdLogger(
		logging,
		func(_ logger.ILogArg) string {
			return ""
		},
	)

	sett := net_message.NewSettings(&net_message.SSettings{
		FWorkSizeBits: 1,
		FNetworkKey:   "_",
	})

	producer := &tsProducer{}
	addr := testutils.TgAddrs[44]
	go func() {
		err := ProduceProcessor(ctx, producer, httpLogger, sett, addr)
		if err != nil {
			return
		}
	}()
	time.Sleep(200 * time.Millisecond)

	msg := net_message.NewMessage(
		net_message.NewConstructSettings(&net_message.SConstructSettings{
			FSettings: sett,
			FParallel: 1,
		}),
		payload.NewPayload32(1, []byte("hello, world!")),
	)

	urlAddr := "http://" + addr + "/adapter"
	client := &http.Client{Timeout: time.Minute}
	if _, err := api.Request(ctx, client, http.MethodPost, urlAddr, msg.ToString()); err != nil {
		t.Error(err)
		return
	}

	if _, err := api.Request(ctx, client, http.MethodGet, urlAddr, msg.ToString()); err == nil {
		t.Error("success request with invalid method")
		return
	}
	if _, err := api.Request(ctx, client, http.MethodPost, urlAddr, []byte{1}); err == nil {
		t.Error("success request with invalid message")
		return
	}

	producer.fNeedFailed = true
	if _, err := api.Request(ctx, client, http.MethodPost, urlAddr, msg.ToString()); err == nil {
		t.Error("success request with produce failed")
		return
	}
}

type tsProducer struct {
	fNeedFailed bool
}
type tsConsumer struct {
	fErr chan error
	fMsg chan net_message.IMessage
}

func (p *tsProducer) Produce(_ context.Context, _ net_message.IMessage) error {
	if p.fNeedFailed {
		return errors.New("some error") // nolint: err113
	}
	return nil
}

func (p *tsConsumer) Consume(ctx context.Context) (net_message.IMessage, error) {
	select {
	case <-ctx.Done():
		return nil, nil
	case x := <-p.fMsg:
		return x, nil
	case x := <-p.fErr:
		return nil, x
	}
}

type tsHLTClient struct {
	fNeedFailed bool
}

func (p *tsHLTClient) GetIndex(context.Context) (string, error) { return "", nil }
func (p *tsHLTClient) GetSettings(context.Context) (hlt_config.IConfigSettings, error) {
	return nil, nil
}

func (p *tsHLTClient) GetPointer(context.Context) (uint64, error)      { return 0, nil }
func (p *tsHLTClient) GetHash(context.Context, uint64) (string, error) { return "", nil }

func (p *tsHLTClient) GetMessage(context.Context, string) (net_message.IMessage, error) {
	return nil, nil
}
func (p *tsHLTClient) PutMessage(context.Context, net_message.IMessage) error {
	if p.fNeedFailed {
		return errors.New("some error") // nolint: err113
	}
	return nil
}
