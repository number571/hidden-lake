package utils

import (
	"context"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
	hls_config "github.com/number571/hidden-lake/internal/service/pkg/config"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SUtilsError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestUtils(t *testing.T) {
	t.Parallel()

	limit, err := GetMessageLimit(context.Background(), newTsHLSClient(1024))
	if err != nil {
		t.Error(err)
		return
	}
	if limit != 1024-gRespSize {
		t.Error("limit != 1024-gRespSize")
		return
	}

	if _, err := GetMessageLimit(context.Background(), newTsHLSClient(gRespSize)); err == nil {
		t.Error("success get message limit with gRespSize >= limit")
		return
	}
}

var (
	_ hls_client.IClient = &tsHLSClient{}
)

type tsHLSClient struct {
	fMsgSize uint64
	fPrivKey asymmetric.IPrivKey
}

func newTsHLSClient(pMsgSize uint64) *tsHLSClient {
	return &tsHLSClient{
		fMsgSize: pMsgSize,
		fPrivKey: asymmetric.NewPrivKey(),
	}
}

func (p *tsHLSClient) GetIndex(context.Context) (string, error) { return "", nil }
func (p *tsHLSClient) GetSettings(context.Context) (hls_config.IConfigSettings, error) {
	return &hls_config.SConfigSettings{
		FPayloadSizeBytes: p.fMsgSize,
	}, nil
}

func (p *tsHLSClient) GetPubKey(context.Context) (asymmetric.IPubKey, error) {
	return p.fPrivKey.GetPubKey(), nil
}

func (p *tsHLSClient) GetOnlines(context.Context) ([]string, error) { return nil, nil }
func (p *tsHLSClient) DelOnline(context.Context, string) error      { return nil }

func (p *tsHLSClient) GetFriends(context.Context) (map[string]asymmetric.IPubKey, error) {
	return nil, nil
}

func (p *tsHLSClient) AddFriend(context.Context, string, asymmetric.IPubKey) error { return nil }
func (p *tsHLSClient) DelFriend(context.Context, string) error                     { return nil }

func (p *tsHLSClient) GetConnections(context.Context) ([]string, error) { return nil, nil }
func (p *tsHLSClient) AddConnection(context.Context, string) error      { return nil }
func (p *tsHLSClient) DelConnection(context.Context, string) error      { return nil }

func (p *tsHLSClient) SendRequest(context.Context, string, request.IRequest) error {
	return nil
}

func (p *tsHLSClient) FetchRequest(context.Context, string, request.IRequest) (response.IResponse, error) {
	return nil, nil
}
