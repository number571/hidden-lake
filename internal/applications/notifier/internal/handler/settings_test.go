// nolint: goerr113
package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/notifier/internal/database"
	"github.com/number571/hidden-lake/internal/applications/notifier/pkg/app/config"
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
	hls_config "github.com/number571/hidden-lake/internal/service/pkg/config"
	"github.com/number571/hidden-lake/internal/utils/language"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"
)

func TestSettingsPage(t *testing.T) {
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
	cfg := &config.SConfig{
		FSettings: &config.SConfigSettings{
			FLanguage: "ENG",
		},
	}

	handler := SettingsPage(ctx, httpLogger, &tsWrapper{fCfg: cfg}, newTsHLSClient(true, true))

	if err := settingsRequestPutOK(handler); err != nil {
		t.Error(err)
		return
	}
	if err := settingsRequestPostOK(handler); err != nil {
		t.Error(err)
		return
	}
	if err := settingsRequestDeleteOK(handler); err != nil {
		t.Error(err)
		return
	}

	if err := settingsRequestDeleteAddress(handler); err == nil {
		t.Error("request success with invalid address")
		return
	}
	if err := settingsRequestPostHostPort(handler); err == nil {
		t.Error("request success with invalid host:port")
		return
	}
	if err := settingsRequestPutLanguage(handler); err == nil {
		t.Error("request success with invalid language")
		return
	}
	if err := settingsRequest404(handler); err == nil {
		t.Error("request success with invalid path")
		return
	}

	handlerx := SettingsPage(ctx, httpLogger, &tsWrapper{fCfg: cfg, fWithFail: true}, newTsHLSClient(true, true))
	if err := settingsRequestPutOK(handlerx); err == nil {
		t.Error("success update language with invalid update config")
		return
	}

	handlery := SettingsPage(ctx, httpLogger, &tsWrapper{fCfg: cfg}, newTsHLSClient(true, false))
	if err := settingsRequestPostOK(handlery); err == nil {
		t.Error("success update conns with invalid add_connection")
		return
	}
	if err := settingsRequestDeleteOK(handlery); err == nil {
		t.Error("success update conns with invalid del_connection")
		return
	}

	handler1 := SettingsPage(ctx, httpLogger, &tsWrapper{fCfg: cfg}, newTsHLSClient(true, false))
	if err := settingsRequestOK(handler1); err == nil {
		t.Error("success get settings with invalid get_pub_key")
		return
	}
}

func settingsRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/settings", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

func settingsRequestDeleteOK(handler http.HandlerFunc) error {
	formData := url.Values{
		"method":  {"DELETE"},
		"address": {"127.0.0.1:9999"},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/settings", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

func settingsRequestDeleteAddress(handler http.HandlerFunc) error {
	formData := url.Values{
		"method": {"DELETE"},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/settings", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

func settingsRequestPostOK(handler http.HandlerFunc) error {
	formData := url.Values{
		"method": {"POST"},
		"host":   {"127.0.0.1"},
		"port":   {"9999"},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/settings", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

func settingsRequestPostHostPort(handler http.HandlerFunc) error {
	formData := url.Values{
		"method": {"POST"},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/settings", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

func settingsRequestPutOK(handler http.HandlerFunc) error {
	formData := url.Values{
		"method":   {"PUT"},
		"language": {"ENG"},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/settings", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

func settingsRequestPutLanguage(handler http.HandlerFunc) error {
	formData := url.Values{
		"method":   {"PUT"},
		"language": {"unknown"},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/settings", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

func settingsRequest404(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/settings/undefined", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

type tsWrapper struct {
	fCfg      config.IConfig
	fWithFail bool
}

func (p *tsWrapper) GetConfig() config.IConfig { return p.fCfg }
func (p *tsWrapper) GetEditor() config.IEditor {
	return &tsEditor{
		fWithFail: p.fWithFail,
	}
}

type tsEditor struct {
	fWithFail bool
}

func (p *tsEditor) UpdateLanguage(language.ILanguage) error {
	if p.fWithFail {
		return errors.New("some error") // nolint: err113
	}
	return nil
}

func (p *tsEditor) UpdateChannels([]string) error {
	if p.fWithFail {
		return errors.New("some error") // nolint: err113
	}
	return nil
}

var (
	_ hls_client.IClient = &tsHLSClient{}
)

type tsHLSClient struct {
	fWithOK       bool
	fGetPubKey    bool
	fPrivKey      asymmetric.IPrivKey
	fFriendPubKey asymmetric.IPubKey
	fPldSize      uint64
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
		FPayloadSizeBytes: 1024,
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
	fMsg    database.IMessage
}

func newTsDatabase(pPushOK, pLoadOK bool) *tsDatabase {
	return &tsDatabase{
		fPushOK: pPushOK,
		fLoadOK: pLoadOK,
	}
}

func (p *tsDatabase) Close() error { return nil }

func (p *tsDatabase) SetHash(_ asymmetric.IPubKey, _ bool, _ []byte) (bool, error) {
	return false, nil
}

func (p *tsDatabase) Size(database.IRelation) uint64 {
	if p.fMsg == nil {
		return 0
	}
	return 1
}

func (p *tsDatabase) Push(_ database.IRelation, pM database.IMessage) error {
	if !p.fPushOK {
		return errors.New("some error") // nolint: err113
	}
	p.fMsg = pM
	return nil
}

func (p *tsDatabase) Load(database.IRelation, uint64, uint64) ([]database.IMessage, error) {
	if !p.fLoadOK {
		return nil, errors.New("some error") // nolint: err113
	}
	if p.fMsg == nil {
		return nil, nil
	}
	return []database.IMessage{p.fMsg}, nil
}
