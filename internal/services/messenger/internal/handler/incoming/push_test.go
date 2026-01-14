package incoming

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	hls_config "github.com/number571/hidden-lake/internal/kernel/pkg/config"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	"github.com/number571/hidden-lake/internal/services/messenger/internal/database"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/message"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"
)

func TestHandleIncomingPushHTTP(t *testing.T) {
	t.Parallel()

	logging, err := std_logger.LoadLogging([]string{})
	if err != nil {
		t.Fatal(err)
	}

	httpLogger := std_logger.NewStdLogger(
		logging,
		func(_ logger.ILogArg) string {
			return ""
		},
	)

	ctx := context.Background()
	msgBroker := message.NewMessageBroker()
	handler := HandleIncomingPushHTTP(ctx, httpLogger, newTsDatabase(true, true), msgBroker, newTsHLSClient(true, true))

	if err := incomingPushRequestOK(handler); err != nil {
		t.Fatal(err)
	}

	if err := incomingPushRequestMethod(handler); err == nil {
		t.Fatal("request success with invalid method")
	}
	if err := incomingPushRequestPubKey(handler); err == nil {
		t.Fatal("request success with invalid pubkey")
	}
	if err := incomingPushRequestMessage(handler); err == nil {
		t.Fatal("request success with invalid message")
	}

	handlerx := HandleIncomingPushHTTP(ctx, httpLogger, newTsDatabase(true, true), msgBroker, newTsHLSClient(false, true))
	if err := incomingPushRequestOK(handlerx); err == nil {
		t.Fatal("request success with invalid my pubkey")
	}
	handlery := HandleIncomingPushHTTP(ctx, httpLogger, newTsDatabase(false, true), msgBroker, newTsHLSClient(true, true))
	if err := incomingPushRequestOK(handlery); err == nil {
		t.Fatal("request success with invalid push message")
	}
}

func incomingPushRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/push", bytes.NewBuffer([]byte("hello, world!")))
	req.Header.Set(hlk_settings.CHeaderSenderName, "abc")

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

func incomingPushRequestMessage(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/push", nil)
	req.Header.Set(hlk_settings.CHeaderSenderName, "abc")

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

func incomingPushRequestPubKey(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/push", bytes.NewBuffer([]byte("hello, world!")))

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

func incomingPushRequestMethod(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/push", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

var (
	_ hlk_client.IClient = &tsHLSClient{}
)

type tsHLSClient struct {
	fWithOK       bool
	fGetPubKey    bool
	fPrivKey      asymmetric.IPrivKey
	fFriendPubKey asymmetric.IPubKey
}

func newTsHLSClient(pGetPubKey, pWithOK bool) *tsHLSClient {
	return &tsHLSClient{
		fWithOK:       pWithOK,
		fGetPubKey:    pGetPubKey,
		fPrivKey:      asymmetric.NewPrivKey(),
		fFriendPubKey: asymmetric.NewPrivKey().GetPubKey(),
	}
}

func (p *tsHLSClient) GetIndex(context.Context) (string, error) { return "", nil }
func (p *tsHLSClient) GetSettings(context.Context) (hls_config.IConfigSettings, error) {
	if !p.fWithOK {
		return nil, errors.New("some error") // nolint: err113
	}
	return &hls_config.SConfigSettings{
		FPayloadSizeBytes: 256,
	}, nil
}

func (p *tsHLSClient) GetPubKey(context.Context) (asymmetric.IPubKey, error) {
	if !p.fGetPubKey {
		return nil, errors.New("some error") // nolint: err113
	}
	return p.fPrivKey.GetPubKey(), nil
}

func (p *tsHLSClient) GetOnlines(context.Context) ([]string, error) {
	return []string{"tcp://aaa"}, nil
}
func (p *tsHLSClient) DelOnline(context.Context, string) error { return nil }

func (p *tsHLSClient) GetFriends(context.Context) (map[string]asymmetric.IPubKey, error) {
	return map[string]asymmetric.IPubKey{
		"abc": p.fFriendPubKey,
	}, nil
}

func (p *tsHLSClient) AddFriend(context.Context, string, asymmetric.IPubKey) error { return nil }
func (p *tsHLSClient) DelFriend(context.Context, string) error                     { return nil }

func (p *tsHLSClient) GetConnections(context.Context) ([]string, error) {
	return []string{"tcp://aaa"}, nil
}
func (p *tsHLSClient) AddConnection(context.Context, string) error {
	if !p.fWithOK {
		return errors.New("some error") // nolint: err113
	}
	return nil
}
func (p *tsHLSClient) DelConnection(context.Context, string) error {
	if !p.fWithOK {
		return errors.New("some error") // nolint: err113
	}
	return nil
}

func (p *tsHLSClient) SendRequest(context.Context, string, request.IRequest) error {
	return nil
}

func (p *tsHLSClient) FetchRequest(context.Context, string, request.IRequest) (response.IResponse, error) {
	return response.NewResponseBuilder().WithCode(200).Build(), nil
}

type tsDatabase struct {
	fPushOK bool
	fLoadOK bool
	fMsg    message.IMessage
}

func newTsDatabase(pPushOK, pLoadOK bool) *tsDatabase {
	return &tsDatabase{
		fPushOK: pPushOK,
		fLoadOK: pLoadOK,
	}
}

func (p *tsDatabase) Close() error { return nil }

func (p *tsDatabase) Size(database.IRelation) uint64 {
	if p.fMsg == nil {
		return 0
	}
	return 1
}

func (p *tsDatabase) Push(_ database.IRelation, pM message.IMessage) error {
	if !p.fPushOK {
		return errors.New("some error") // nolint: err113
	}
	p.fMsg = pM
	return nil
}

func (p *tsDatabase) Load(database.IRelation, uint64, uint64) ([]message.IMessage, error) {
	if !p.fLoadOK {
		return nil, errors.New("some error") // nolint: err113
	}
	if p.fMsg == nil {
		return nil, nil
	}
	return []message.IMessage{p.fMsg}, nil
}
