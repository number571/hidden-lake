package handler

// import (
// 	"context"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/number571/go-peer/pkg/crypto/asymmetric"
// 	"github.com/number571/go-peer/pkg/logger"
// 	"github.com/number571/hidden-lake/internal/applications/notifier/pkg/app/config"
// 	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
// 	hls_config "github.com/number571/hidden-lake/internal/service/pkg/config"
// 	"github.com/number571/hidden-lake/internal/utils/language"
// 	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
// 	"github.com/number571/hidden-lake/pkg/request"
// 	"github.com/number571/hidden-lake/pkg/response"
// )

// func TestHandleIncomingNotifyHTTP(t *testing.T) {
// 	t.Parallel()

// 	ctx := context.Background()
// 	log := logger.NewLogger(
// 		logger.NewSettings(&logger.SSettings{}),
// 		func(_ logger.ILogArg) string { return "" },
// 	)
// 	handler := HandleIncomingRedirectHTTP(ctx, &tsConfig{}, log, newTsHLSClient(true, true))

// 	if err := incomingRequestMethod(handler); err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if err := incomingRequestSuccess(handler); err != nil {
// 		t.Error(err)
// 		return
// 	}
// }

// func incomingRequestMethod(handler http.HandlerFunc) error {
// 	w := httptest.NewRecorder()
// 	req := httptest.NewRequest(http.MethodPost, "/", nil)

// 	handler(w, req)
// 	res := w.Result()
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusMethodNotAllowed {
// 		return errors.New("bad status code") // nolint: err113
// 	}

// 	return nil
// }

// func incomingRequestSuccess(handler http.HandlerFunc) error {
// 	w := httptest.NewRecorder()
// 	req := httptest.NewRequest(http.MethodGet, "/", nil)

// 	handler(w, req)
// 	res := w.Result()
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		return errors.New("bad status code") // nolint: err113
// 	}

// 	return nil
// }

// type tsConfig struct{}

// func (p *tsConfig) GetSettings() config.IConfigSettings { return &tsConfigSettings{} }
// func (p *tsConfig) GetAddress() config.IAddress         { return nil }
// func (p *tsConfig) GetLogging() std_logger.ILogging     { return nil }
// func (p *tsConfig) GetConnection() string               { return "" }

// type tsConfigSettings struct{}

// func (p *tsConfigSettings) GetMessagesCapacity() uint64 { return 2048 }
// func (p *tsConfigSettings) GetPowParallel() uint64      { return 1 }
// func (p *tsConfigSettings) GetWorkSizeBits() uint64     { return 1 }
// func (p *tsConfigSettings) GetLanguage() language.ILanguage {
// 	return language.ILanguage(language.CLangENG)
// }

// func (p *tsConfigSettings) GetResponseMessage() string { return "" }

// var (
// 	_ hls_client.IClient = &tsHLSClient{}
// )

// type tsHLSClient struct {
// 	fFetchOK bool
// 	fWithOK  bool
// 	fPrivKey asymmetric.IPrivKey
// 	// fPldSize uint64
// }

// func newTsHLSClient(pFetchOK, pWithOK bool) *tsHLSClient { // nolint: unparam
// 	return &tsHLSClient{
// 		fFetchOK: pFetchOK,
// 		fWithOK:  pWithOK,
// 		fPrivKey: asymmetric.NewPrivKey(),
// 	}
// }

// func (p *tsHLSClient) GetIndex(context.Context) (string, error) { return "", nil }
// func (p *tsHLSClient) GetSettings(context.Context) (hls_config.IConfigSettings, error) {
// 	if !p.fWithOK {
// 		return nil, errors.New("some error") // nolint: err113
// 	}
// 	return &hls_config.SConfigSettings{
// 		FPayloadSizeBytes: 1024,
// 	}, nil
// }

// func (p *tsHLSClient) GetPubKey(context.Context) (asymmetric.IPubKey, error) {
// 	return p.fPrivKey.GetPubKey(), nil
// }

// func (p *tsHLSClient) GetOnlines(context.Context) ([]string, error) {
// 	return []string{"tcp://aaa"}, nil
// }
// func (p *tsHLSClient) DelOnline(context.Context, string) error { return nil }

// func (p *tsHLSClient) GetFriends(context.Context) (map[string]asymmetric.IPubKey, error) {
// 	return nil, nil
// }

// func (p *tsHLSClient) AddFriend(context.Context, string, asymmetric.IPubKey) error { return nil }
// func (p *tsHLSClient) DelFriend(context.Context, string) error                     { return nil }

// func (p *tsHLSClient) GetConnections(context.Context) ([]string, error) {
// 	return []string{"tcp://aaa"}, nil
// }
// func (p *tsHLSClient) AddConnection(context.Context, string) error { return nil }
// func (p *tsHLSClient) DelConnection(context.Context, string) error { return nil }

// func (p *tsHLSClient) SendRequest(context.Context, string, request.IRequest) error {
// 	return nil
// }

// func (p *tsHLSClient) FetchRequest(context.Context, string, request.IRequest) (response.IResponse, error) {
// 	if !p.fFetchOK {
// 		return nil, errors.New("some error") // nolint: err113
// 	}
// 	resp := response.NewResponseBuilder().WithCode(200).WithBody([]byte(`[{"name":"file.txt","hash":"114a856f792c4c292599dba6fa41adba45ef4f851b1d17707e2729651968ff64be375af9cff6f9547b878d5c73c16a11","size":500}]`))
// 	return resp.Build(), nil
// }
