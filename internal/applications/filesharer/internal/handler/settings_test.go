// nolint: err113
package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/filesharer/pkg/app/config"
	"github.com/number571/hidden-lake/internal/utils/language"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func TestSettingsPage(t *testing.T) {
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
	cfg := &config.SConfig{
		FSettings: &config.SConfigSettings{
			FLanguage: "ENG",
		},
	}

	handler := SettingsPage(ctx, httpLogger, &tsWrapper{fCfg: cfg}, newTsHLSClient(true, true))

	if err := settingsRequestPutOK(handler); err != nil {
		t.Fatal(err)
	}
	if err := settingsRequestPostOK(handler); err != nil {
		t.Fatal(err)
	}
	if err := settingsRequestDeleteOK(handler); err != nil {
		t.Fatal(err)
	}

	if err := settingsRequestDeleteAddress(handler); err == nil {
		t.Fatal("request success with invalid address")
	}
	if err := settingsRequestPostHostPort(handler); err == nil {
		t.Fatal("request success with invalid host:port")
	}
	if err := settingsRequestPutLanguage(handler); err == nil {
		t.Fatal("request success with invalid language")
	}
	if err := settingsRequest404(handler); err == nil {
		t.Fatal("request success with invalid path")
	}

	handlerx := SettingsPage(ctx, httpLogger, &tsWrapper{fCfg: cfg, fWithFail: true}, newTsHLSClient(true, true))
	if err := settingsRequestPutOK(handlerx); err == nil {
		t.Fatal("success update language with invalid update config")
	}

	handlery := SettingsPage(ctx, httpLogger, &tsWrapper{fCfg: cfg}, newTsHLSClient(true, false))
	if err := settingsRequestPostOK(handlery); err == nil {
		t.Fatal("success update conns with invalid add_connection")
	}
	if err := settingsRequestDeleteOK(handlery); err == nil {
		t.Fatal("success update conns with invalid del_connection")
	}

	handler1 := SettingsPage(ctx, httpLogger, &tsWrapper{fCfg: cfg}, newTsHLSClient(true, false))
	if err := settingsRequestOK(handler1); err == nil {
		t.Fatal("success get settings with invalid get_pub_key")
	}
}

func settingsRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/settings", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
func (p *tsWrapper) GetEditor() config.IEditor { return &tsEditor{p.fWithFail} }

type tsEditor struct {
	fWithFail bool
}

func (p *tsEditor) UpdateLanguage(language.ILanguage) error {
	if p.fWithFail {
		return errors.New("some error") // nolint: err113
	}
	return nil
}
