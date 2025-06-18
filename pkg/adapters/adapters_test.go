package adapters

import (
	"context"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/hidden-lake/build"
)

func TestSettings(t *testing.T) {
	t.Parallel()

	sett := NewSettings(nil)
	defaultNetwork, _ := build.GetNetwork(build.CDefaultNetwork)

	if sett.GetMessageSizeBytes() != defaultNetwork.FMessageSizeBytes {
		t.Fatal("get invalid settings")
	}

	_ = NewSettingsByNetworkKey(build.CDefaultNetwork)
}

func TestPanicSettings(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("nothing panics")
		}
	}()

	_ = NewSettingsByNetworkKey("__test_unknown__")
}

func TestNewRunnerAdapter(t *testing.T) {
	t.Parallel()

	adapter := NewRunnerAdapter(&tsAdapter{}, func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() { _ = adapter.Run(ctx) }()

	time.Sleep(100 * time.Millisecond)
	cancel()
}

type tsAdapter struct{}

func (p *tsAdapter) Consume(context.Context) (layer1.IMessage, error) { return nil, nil }
func (p *tsAdapter) Produce(context.Context, layer1.IMessage) error   { return nil }
