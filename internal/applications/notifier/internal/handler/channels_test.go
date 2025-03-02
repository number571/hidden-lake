// nolint: goerr113
package handler

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/notifier/pkg/app/config"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func TestChannelsPage(t *testing.T) {
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

	handler := ChannelsPage(httpLogger, &tsWrapper{fCfg: &config.SConfig{
		FSettings: &config.SConfigSettings{
			FLanguage: "ENG",
		},
	}})

	if err := channelsRequestOK(handler); err != nil {
		t.Error(err)
		return
	}
	if err := channelsRequestPostOK(handler, "123"); err != nil {
		t.Error(err)
		return
	}
	if err := channelsDeletePostOK(handler, "123"); err != nil {
		t.Error(err)
		return
	}

	if err := channelsRequestPostOK(handler, ""); err == nil {
		t.Error("request success with invalid key (create)")
		return
	}
	if err := channelsDeletePostOK(handler, ""); err == nil {
		t.Error("request success with invalid key (delete)")
		return
	}
	if err := channelsRequest404(handler); err == nil {
		t.Error("request success with invalid path")
		return
	}
}

func channelsRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/channels", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func channelsRequest404(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/channels/undefined", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

func channelsRequestPostOK(handler http.HandlerFunc, key string) error {
	formData := url.Values{
		"method": {"POST"},
		"key":    {key},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/channels", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

func channelsDeletePostOK(handler http.HandlerFunc, key string) error {
	formData := url.Values{
		"method": {"DELETE"},
		"key":    {key},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/channels", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}
