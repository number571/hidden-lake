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

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/messenger/internal/database"
	"github.com/number571/hidden-lake/internal/applications/messenger/pkg/app/config"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func TestFriendsChatPage(t *testing.T) {
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

	db := newTsDatabase(true)
	_ = db.Push(nil, database.NewMessage(false, wrapText("hello, world!")))

	handler := FriendsChatPage(ctx, httpLogger, cfg, db, newTsHLSClient(true))
	if err := friendsChatRequestOK(handler); err != nil {
		t.Error(err)
		return
	}
	if err := friendsChatRequestPostOK(handler); err != nil {
		t.Error(err)
		return
	}
	if err := friendsChatRequestPostPingOK(handler); err != nil {
		t.Error(err)
		return
	}

	if err := friendsChatRequest404(handler); err == nil {
		t.Error("request success with invalid path")
		return
	}
}

func friendsChatRequestPostOK(handler http.HandlerFunc) error {
	formData := url.Values{
		"method":        {"POST"},
		"input_message": {"hello, world!"},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/friends/chat?alias_name=abc", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusSeeOther {
		return errors.New("bad status code")
	}

	return nil
}

func friendsChatRequestPostPingOK(handler http.HandlerFunc) error {
	formData := url.Values{
		"method": {"POST"},
		"ping":   {"true"},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/friends/chat?alias_name=abc", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

func friendsChatRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/friends/chat?alias_name=abc", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

func friendsChatRequest404(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/friends/chat/undefined", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}
