package limiters

import (
	"context"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	"github.com/number571/hidden-lake/pkg/api/kernel/client/scheme"
	hls_config "github.com/number571/hidden-lake/pkg/api/kernel/config"
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

func TestUtils(t *testing.T) {
	t.Parallel()

	limit, err := GetLimitOnLoadResponseSize(context.Background(), newTsHLSClient(1024))
	if err != nil {
		t.Fatal(err)
	}
	if limit != 1024-gLoadRspSize {
		t.Fatal("limit != 1024-gRespSize")
	}

	if _, err := GetLimitOnLoadResponseSize(context.Background(), newTsHLSClient(gLoadRspSize)); err == nil {
		t.Fatal("success get message limit with gRespSize >= limit")
	}
}

var (
	_ hlk_client.IClient = &tsHLSClient{}
)

type tsHLSClient struct {
	fMsgSize uint64
}

func newTsHLSClient(pMsgSize uint64) *tsHLSClient {
	return &tsHLSClient{
		fMsgSize: pMsgSize,
	}
}

func (p *tsHLSClient) GetIndex(context.Context) error { return nil }
func (p *tsHLSClient) GetSettings(context.Context) (hls_config.IConfigSettings, error) {
	return &hls_config.SConfigSettings{
		FPayloadSizeBytes: p.fMsgSize,
	}, nil
}

func (p *tsHLSClient) GetOnlines(context.Context) ([]string, error) { return nil, nil }
func (p *tsHLSClient) DelOnline(context.Context, string) error      { return nil }

func (p *tsHLSClient) GetFriends(context.Context, scheme.ISchemeType) (map[string]layer2.IParticipantKey, error) {
	return nil, nil
}

func (p *tsHLSClient) AddFriend(context.Context, string, layer2.IParticipantKey) error { return nil }
func (p *tsHLSClient) DelFriend(context.Context, string) error                         { return nil }

func (p *tsHLSClient) GetConnections(context.Context) ([]string, error) { return nil, nil }
func (p *tsHLSClient) AddConnection(context.Context, string) error      { return nil }
func (p *tsHLSClient) DelConnection(context.Context, string) error      { return nil }

func (p *tsHLSClient) SendRequest(context.Context, string, request.IRequest) error {
	return nil
}

func (p *tsHLSClient) FetchRequest(context.Context, string, request.IRequest) (response.IResponse, error) {
	return nil, nil
}
