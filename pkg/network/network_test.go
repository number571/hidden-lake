package network

import (
	"context"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"
)

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
			FMessageSizeBytes: 4570,
			FQueuePeriod:      time.Second,
			FFetchTimeout:     time.Second,
		}),
		asymmetric.NewPrivKey(),
		&tsDatabase{},
		func() []string { return nil },
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
