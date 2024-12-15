package network

import (
	"context"
	"testing"
	"time"

	gopeer_adapters "github.com/number571/go-peer/pkg/anonymity/adapters"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
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
			FQueuePeriod:  time.Second,
			FFetchTimeout: time.Second,
		}),
		asymmetric.NewPrivKey(),
		&tsDatabase{},
		adapters.NewRunnerAdapter(
			gopeer_adapters.NewAdapterByFuncs(
				func(context.Context, net_message.IMessage) error { return nil },
				func(context.Context) (net_message.IMessage, error) { return nil, nil },
			),
			func(context.Context) error { return nil },
		),
		func(_ context.Context, _ asymmetric.IPubKey, _ request.IRequest) (response.IResponse, error) {
			return nil, nil
		},
	)
}

func TestPanicSettings(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	_ = NewSettingsByNetworkKey("__test_unknown__", nil)
}

func TestSettings(t *testing.T) {
	t.Parallel()

	_ = NewSettingsByNetworkKey(build.CDefaultNetwork, nil)
}

type tsDatabase struct{}

func (p *tsDatabase) Close() error               { return nil }
func (p *tsDatabase) Set([]byte, []byte) error   { return nil }
func (p *tsDatabase) Get([]byte) ([]byte, error) { return nil, nil }
func (p *tsDatabase) Del([]byte) error           { return nil }
