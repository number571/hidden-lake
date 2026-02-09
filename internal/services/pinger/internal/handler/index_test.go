package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	hlk_config "github.com/number571/hidden-lake/pkg/api/kernel/config"
	hls_client "github.com/number571/hidden-lake/pkg/api/services/pinger/client"
	"github.com/number571/hidden-lake/pkg/network/request"
	"github.com/number571/hidden-lake/pkg/network/response"
)

func TestErrorsAPI(t *testing.T) {
	t.Parallel()

	client := hls_client.NewClient(
		hls_client.NewRequester("", &http.Client{}),
	)

	if err := client.GetIndex(context.Background()); err == nil {
		t.Fatal("success incorrect getIndex")
	}
	if err := client.PingFriend(context.Background(), ""); err == nil {
		t.Fatal("success incorrect pingFriend")
	}
}

func TestHandleIndexAPI(t *testing.T) {
	t.Parallel()

	httpLogger := std_logger.NewStdLogger(
		func() std_logger.ILogging {
			logging, err := std_logger.LoadLogging([]string{})
			if err != nil {
				panic(err)
			}
			return logging
		}(),
		func(_ logger.ILogArg) string {
			return ""
		},
	)

	handler := HandleIndexAPI(httpLogger)
	if err := indexAPIRequestOK(handler); err != nil {
		t.Fatal(err)
	}
}

func indexAPIRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

var (
	_ hlk_client.IClient = &tsHLKClient{}
)

type tsHLKClient struct {
	fFetchType int
	fPrivKey   asymmetric.IPrivKey
}

func newTsHLKClient(pFetchType int) *tsHLKClient {
	return &tsHLKClient{
		fFetchType: pFetchType,
		fPrivKey:   asymmetric.NewPrivKey(),
	}
}

func (p *tsHLKClient) GetIndex(context.Context) error { return nil }
func (p *tsHLKClient) GetSettings(context.Context) (hlk_config.IConfigSettings, error) {
	return &hlk_config.SConfigSettings{
		FPayloadSizeBytes: 1024,
	}, nil
}

func (p *tsHLKClient) GetPubKey(context.Context) (asymmetric.IPubKey, error) {
	return p.fPrivKey.GetPubKey(), nil
}

func (p *tsHLKClient) GetOnlines(context.Context) ([]string, error) {
	return []string{"tcp://aaa"}, nil
}
func (p *tsHLKClient) DelOnline(context.Context, string) error { return nil }

func (p *tsHLKClient) GetFriends(context.Context) (map[string]asymmetric.IPubKey, error) {
	return map[string]asymmetric.IPubKey{
		"abc": asymmetric.NewPrivKey().GetPubKey(),
	}, nil
}

func (p *tsHLKClient) AddFriend(context.Context, string, asymmetric.IPubKey) error { return nil }
func (p *tsHLKClient) DelFriend(context.Context, string) error                     { return nil }

func (p *tsHLKClient) GetConnections(context.Context) ([]string, error) {
	return []string{"tcp://aaa"}, nil
}
func (p *tsHLKClient) AddConnection(context.Context, string) error { return nil }
func (p *tsHLKClient) DelConnection(context.Context, string) error { return nil }

func (p *tsHLKClient) SendRequest(context.Context, string, request.IRequest) error {
	return nil
}

func (p *tsHLKClient) FetchRequest(context.Context, string, request.IRequest) (response.IResponse, error) {
	switch p.fFetchType {
	case 1:
		return nil, errors.New("error") // nolint: err113
	case 0:
		resp := response.NewResponseBuilder().WithCode(200).WithBody([]byte(`pong`))
		return resp.Build(), nil
	case -1:
		resp := response.NewResponseBuilder().WithCode(500).WithBody([]byte(`500`))
		return resp.Build(), nil
	case -2:
		resp := response.NewResponseBuilder().WithCode(200).WithBody([]byte{1})
		return resp.Build(), nil
	}
	panic("unknown fetch type")
}
