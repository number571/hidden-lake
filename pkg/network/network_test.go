package network

import (
	"context"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	hiddenlake "github.com/number571/hidden-lake"
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

	for i := 0; i < 4; i++ {
		testPanicSettings(t, i)
	}
}

func testPanicSettings(t *testing.T, n int) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	switch n {
	case 0:
		_ = NewSettingsByNetworkKey("__test_unknown__", nil)
	case 1:
		_ = NewSettings(&SSettings{
			FMessageSizeBytes: 8192,
			FQueuePeriod:      time.Second,
		})
	case 2:
		_ = NewSettings(&SSettings{
			FMessageSizeBytes: 8192,
			FFetchTimeout:     time.Second,
		})
	case 3:
		_ = NewSettings(&SSettings{
			FQueuePeriod:  time.Second,
			FFetchTimeout: time.Second,
		})
	}
}

func TestSettings(t *testing.T) {
	t.Parallel()

	_ = NewSettingsByNetworkKey(hiddenlake.CDefaultNetwork, nil)
}

type tsDatabase struct{}

func (p *tsDatabase) Close() error               { return nil }
func (p *tsDatabase) Set([]byte, []byte) error   { return nil }
func (p *tsDatabase) Get([]byte) ([]byte, error) { return nil, nil }
func (p *tsDatabase) Del([]byte) error           { return nil }
