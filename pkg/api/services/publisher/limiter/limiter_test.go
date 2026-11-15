package limiter

import (
	"context"
	"testing"

	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	hlk_config "github.com/number571/hidden-lake/pkg/api/kernel/config"

	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/hidden-lake/pkg/api/kernel/client/scheme"
	"github.com/number571/hidden-lake/pkg/network/request"
	"github.com/number571/hidden-lake/pkg/network/response"
)

type tsTest struct {
	d string
	f func(context.Context, hlk_client.IClient) (uint64, error)
	r uint64
}

func TestLimiter(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	hlkClient := &tsHLKClient{144} // 200 - 56(HEADER)

	tt := []tsTest{
		{d: "push_request", f: GetLimitOnPushRequestSize, r: 65},
		{d: "push_response", f: GetLimitOnPushResponseSize, r: 104},
		{d: "load_request", f: GetLimitOnLoadRequestSize, r: 50},
		{d: "load_response", f: GetLimitOnLoadResponseSize, r: 74},
		{d: "size_request", f: GetLimitOnSizeRequestSize, r: 83},
		{d: "size_response", f: GetLimitOnSizeResponseSize, r: 104},
	}

	for _, tt := range tt {
		v, err := tt.f(ctx, hlkClient)
		if err != nil {
			t.Fatalf("%s: %s", tt.d, err)
		}
		if v != tt.r {
			t.Fatalf("%s: got(%d) != except(%d)", tt.d, v, tt.r)
		}
	}
}

type tsHLKClient struct {
	fPldSoze uint64
}

func (p *tsHLKClient) GetIndex(context.Context) error { return nil }
func (p *tsHLKClient) GetSettings(context.Context) (hlk_config.IConfigSettings, error) {
	return &hlk_config.SConfigSettings{FPayloadSizeBytes: p.fPldSoze}, nil
}

func (p *tsHLKClient) GetOnlines(context.Context) ([]string, error) {
	return []string{}, nil
}
func (p *tsHLKClient) DelOnline(context.Context, string) error { return nil }

func (p *tsHLKClient) GetFriends(context.Context, scheme.ISchemeType) (map[string]layer2.IParticipantKey, error) {
	return map[string]layer2.IParticipantKey{}, nil
}

func (p *tsHLKClient) AddFriend(context.Context, string, layer2.IParticipantKey) error { return nil }
func (p *tsHLKClient) DelFriend(context.Context, string) error                         { return nil }

func (p *tsHLKClient) GetConnections(context.Context) ([]string, error) {
	return []string{}, nil
}
func (p *tsHLKClient) AddConnection(context.Context, string) error { return nil }
func (p *tsHLKClient) DelConnection(context.Context, string) error { return nil }

func (p *tsHLKClient) SendRequest(context.Context, string, request.IRequest) error {
	return nil
}

func (p *tsHLKClient) FetchRequest(_ context.Context, _ string, pReq request.IRequest) (response.IResponse, error) {
	return nil, nil
}
