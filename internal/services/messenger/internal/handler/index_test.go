package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/messenger/internal/database"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/app/config"
	logger_std "github.com/number571/hidden-lake/internal/utils/logger/std"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	hlk_config "github.com/number571/hidden-lake/pkg/api/kernel/config"
	hls_client "github.com/number571/hidden-lake/pkg/api/services/messenger/client"
	"github.com/number571/hidden-lake/pkg/api/services/messenger/client/dto"
	"github.com/number571/hidden-lake/pkg/network/request"
	"github.com/number571/hidden-lake/pkg/network/response"
)

func TestErrorsAPI(t *testing.T) {
	t.Parallel()

	client := hls_client.NewClient(
		hls_client.NewRequester("", &http.Client{}),
	)

	if _, err := client.GetIndex(context.Background()); err == nil {
		t.Fatal("success incorrect getIndex")
	}
	if _, err := client.GetMessageLimit(context.Background()); err == nil {
		t.Fatal("success incorrect getMessageLimit")
	}
	if _, err := client.ListenChat(context.Background(), "", ""); err == nil {
		t.Fatal("success incorrect listenChat")
	}
	if _, err := client.PushMessage(context.Background(), "", ""); err == nil {
		t.Fatal("success incorrect pushMessage")
	}
	if _, err := client.LoadMessages(context.Background(), "", 0, 0, false); err == nil {
		t.Fatal("success incorrect loadMessages")
	}
	if _, err := client.CountMessages(context.Background(), ""); err == nil {
		t.Fatal("success incorrect countMessages")
	}
}

func TestHandleIndexAPI(t *testing.T) {
	t.Parallel()

	httpLogger := logger_std.NewStdLogger(
		func() logger_std.ILogging {
			logging, err := logger_std.LoadLogging([]string{})
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
	fPubKeyOK   bool
	fSettingsOK bool
	fSendOK     bool
	fPrivKey    asymmetric.IPrivKey
}

func newTsHLKClient(pPubKeyOK bool, pSettingsOK bool, pSendOK bool) *tsHLKClient {
	return &tsHLKClient{
		fSettingsOK: pSettingsOK,
		fPubKeyOK:   pPubKeyOK,
		fSendOK:     pSendOK,
		fPrivKey:    asymmetric.NewPrivKey(),
	}
}

func (p *tsHLKClient) GetIndex(context.Context) (string, error) { return "", nil }
func (p *tsHLKClient) GetSettings(context.Context) (hlk_config.IConfigSettings, error) {
	if !p.fSettingsOK {
		return nil, errors.New("error") // nolint: err113
	}
	return &hlk_config.SConfigSettings{
		FPayloadSizeBytes: 1024,
	}, nil
}

func (p *tsHLKClient) GetPubKey(context.Context) (asymmetric.IPubKey, error) {
	if !p.fPubKeyOK {
		return nil, errors.New("error") // nolint: err113
	}
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
	if !p.fSendOK {
		return errors.New("error") // nolint: err113
	}
	return nil
}

func (p *tsHLKClient) FetchRequest(context.Context, string, request.IRequest) (response.IResponse, error) {
	resp := response.NewResponseBuilder().WithCode(200).WithBody([]byte{1})
	return resp.Build(), nil
}

type tsConfig struct {
}

func (p *tsConfig) GetSettings() config.IConfigSettings {
	return &config.SConfigSettings{
		FMessagesCapacity: 128,
	}
}
func (p *tsConfig) GetAddress() config.IAddress {
	return nil
}
func (p *tsConfig) GetLogging() logger_std.ILogging {
	return nil
}
func (p *tsConfig) GetConnection() string {
	return ""
}

type tsDatabase struct {
	fLoadOK bool
	fPushOK bool
}

func newTsDatabase(pLoadOK, pPushOK bool) *tsDatabase {
	return &tsDatabase{
		fLoadOK: pLoadOK,
		fPushOK: pPushOK,
	}
}

func (p *tsDatabase) Close() error {
	return nil
}
func (p *tsDatabase) Size(database.IRelation) uint64 {
	return 1
}
func (p *tsDatabase) Push(database.IRelation, dto.IMessage) error {
	if !p.fPushOK {
		return errors.New("error") // nolint: err113
	}
	return nil
}
func (p *tsDatabase) Load(_ database.IRelation, x uint64, o uint64) ([]dto.IMessage, error) {
	if !p.fLoadOK {
		return nil, errors.New("error") // nolint: err113
	}
	if x != 0 || o == 0 {
		return nil, nil
	}
	return []dto.IMessage{
		dto.NewMessage(true, "hello, world!", time.Now()),
	}, nil
}
