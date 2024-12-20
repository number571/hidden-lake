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
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
	hls_config "github.com/number571/hidden-lake/internal/service/pkg/config"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"
)

func TestHandleIncomingLoadHTTP(t *testing.T) {
	t.Parallel()

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

	ctx := context.Background()

	handler := HandleIncomingLoadHTTP(ctx, httpLogger, "./testdata", newTsHLSClient(true, true))
	if err := incomingLoadRequestOK(handler); err != nil {
		t.Error(err)
		return
	}

	if err := incomingLoadRequestMethod(handler); err == nil {
		t.Error("request success with invalid method")
		return
	}
	if err := incomingLoadRequestChunk(handler); err == nil {
		t.Error("request success with invalid chunk")
		return
	}
	if err := incomingLoadRequestName(handler); err == nil {
		t.Error("request success with invalid name")
		return
	}
	if err := incomingLoadRequestFile(handler); err == nil {
		t.Error("request success with invalid file")
		return
	}
	if err := incomingLoadRequestSize(handler); err == nil {
		t.Error("request success with invalid size")
		return
	}
	if err := incomingLoadRequestNotFound(handler); err == nil {
		t.Error("request success with not found file")
		return
	}
	if err := incomingLoadRequestBigChunk(handler); err == nil {
		t.Error("request success with big chunk number")
		return
	}

	handlerx := HandleIncomingLoadHTTP(ctx, httpLogger, "./testdata", newTsHLSClient(true, false))
	if err := incomingLoadRequestOK(handlerx); err == nil {
		t.Error("success request with failed get message size")
		return
	}
}

func incomingLoadRequestBigChunk(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/load?name=file.txt&chunk=10000", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func incomingLoadRequestNotFound(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/load?name=norfound.txt&chunk=0", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func incomingLoadRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/load?name=file.txt&chunk=0", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func incomingLoadRequestSize(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/load?file=file.txt&chunk=1024", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingLoadRequestFile(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/load?file=notfound.txt&chunk=0", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingLoadRequestName(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/load?chunk=0", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingLoadRequestChunk(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/load?name=file.txt", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingLoadRequestMethod(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/load", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

var (
	_ hls_client.IClient = &tsHLSClient{}
)

type tsHLSClient struct {
	fFetchOK bool
	fWithOK  bool
	fPrivKey asymmetric.IPrivKey
}

func newTsHLSClient(pFetchOK, pWithOK bool) *tsHLSClient {
	return &tsHLSClient{
		fFetchOK: pFetchOK,
		fWithOK:  pWithOK,
		fPrivKey: asymmetric.NewPrivKey(),
	}
}

func (p *tsHLSClient) GetIndex(context.Context) (string, error) { return "", nil }
func (p *tsHLSClient) GetSettings(context.Context) (hls_config.IConfigSettings, error) {
	if !p.fWithOK {
		return nil, errors.New("some error") // nolint: err113
	}
	return &hls_config.SConfigSettings{
		FLimitMessageSizeBytes: 1024,
	}, nil
}

func (p *tsHLSClient) GetPubKey(context.Context) (asymmetric.IPubKey, error) {
	return p.fPrivKey.GetPubKey(), nil
}

func (p *tsHLSClient) GetOnlines(context.Context) ([]string, error) {
	return []string{"tcp://aaa"}, nil
}
func (p *tsHLSClient) DelOnline(context.Context, string) error { return nil }

func (p *tsHLSClient) GetFriends(context.Context) (map[string]asymmetric.IPubKey, error) {
	return nil, nil
}

func (p *tsHLSClient) AddFriend(context.Context, string, asymmetric.IPubKey) error { return nil }
func (p *tsHLSClient) DelFriend(context.Context, string) error                     { return nil }

func (p *tsHLSClient) GetConnections(context.Context) ([]string, error) {
	return []string{"tcp://aaa"}, nil
}
func (p *tsHLSClient) AddConnection(context.Context, string) error { return nil }
func (p *tsHLSClient) DelConnection(context.Context, string) error { return nil }

func (p *tsHLSClient) SendRequest(context.Context, string, request.IRequest) error {
	return nil
}

func (p *tsHLSClient) FetchRequest(context.Context, string, request.IRequest) (response.IResponse, error) {
	if !p.fFetchOK {
		return nil, errors.New("some error") // nolint: err113
	}
	resp := response.NewResponseBuilder().WithCode(200).WithBody([]byte(`[{"name":"file.txt","hash":"114a856f792c4c292599dba6fa41adba45ef4f851b1d17707e2729651968ff64be375af9cff6f9547b878d5c73c16a11","size":500}]`))
	return resp.Build(), nil
}
