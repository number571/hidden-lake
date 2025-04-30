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

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/messenger/pkg/app/config"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func TestFriendsPage(t *testing.T) {
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

	handler := FriendsPage(ctx, httpLogger, cfg, newTsHLSClient(true, true))
	if err := friendsRequestPostOK(handler); err != nil {
		t.Error(err)
		return
	}
	if err := friendsRequestDeleteOK(handler); err != nil {
		t.Error(err)
		return
	}

	if err := friendsRequestPostPubKey(handler); err == nil {
		t.Error("request success with invalid pubkey")
		return
	}
	if err := friendsRequestPostPubKeyVoid(handler); err == nil {
		t.Error("request success with invalid pubkey (void)")
		return
	}
	if err := friendsRequestDeleteAliasName(handler); err == nil {
		t.Error("request success with invalid alias_name")
	}
	if err := friendsRequest404(handler); err == nil {
		t.Error("request success with invalid path")
		return
	}
}

func friendsRequestPostOK(handler http.HandlerFunc) error {
	formData := url.Values{
		"method":     {"POST"},
		"public_key": {asymmetric.NewPrivKey().GetPubKey().ToString()},
		"alias_name": {"alias_name"},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/friends", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

func friendsRequestPostPubKey(handler http.HandlerFunc) error {
	formData := url.Values{
		"method":     {"POST"},
		"public_key": {"abc"},
		"alias_name": {"alias_name"},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/friends", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

func friendsRequestPostPubKeyVoid(handler http.HandlerFunc) error {
	formData := url.Values{
		"method":     {"POST"},
		"alias_name": {"alias_name"},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/friends", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

func friendsRequestDeleteOK(handler http.HandlerFunc) error {
	formData := url.Values{
		"method":     {"DELETE"},
		"alias_name": {"alias_name"},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/friends", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

func friendsRequestDeleteAliasName(handler http.HandlerFunc) error {
	formData := url.Values{
		"method":     {"DELETE"},
		"alias_name": {""},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/friends", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

func friendsRequest404(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/friends/undefined", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}
