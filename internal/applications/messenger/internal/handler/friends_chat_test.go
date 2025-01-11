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
	"github.com/number571/hidden-lake/internal/applications/messenger/internal/utils"
	"github.com/number571/hidden-lake/internal/applications/messenger/pkg/app/config"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func TestGetMessage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	if _, err := utils.GetMessageLimit(ctx, newTsHLSClient(true, false)); err == nil {
		t.Error("success get message limit with settings failed")
		return
	}

	hlsClient := newTsHLSClient(true, false)
	hlsClient.fPldSize = 10
	if _, err := utils.GetMessageLimit(ctx, hlsClient); err == nil {
		t.Error("success get message limit with payload limit")
		return
	}
}

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

	db := newTsDatabase(true, true)
	_ = db.Push(nil, database.NewMessage(false, wrapText("hello, world!")))

	handler := FriendsChatPage(ctx, httpLogger, cfg, db, newTsHLSClient(true, true))
	if err := friendsChatRequestOK(handler); err != nil {
		t.Error(err)
		return
	}
	if err := friendsChatRequestPostOK(handler); err != nil {
		t.Error(err)
		return
	}
	if err := friendsChatRequestPutOK(handler); err != nil {
		t.Error(err)
		return
	}
	if err := friendsChatRequestPostPingOK(handler); err != nil {
		t.Error(err)
		return
	}

	if err := friendsChatRequestPutSizeFile(handler); err == nil {
		t.Error("request success with invalid file size")
		return
	}
	if err := friendsChatRequestHasNotGraphicChars(handler); err == nil {
		t.Error("request success with invalid (has not grpahic chars)")
		return
	}
	if err := friendsChatRequestInputMessage(handler); err == nil {
		t.Error("request success with invalid input_message")
		return
	}
	if err := friendsChatRequestAliasNameNotFound(handler); err == nil {
		t.Error("request success with not found alias_name")
		return
	}
	if err := friendsChatRequestAliasName(handler); err == nil {
		t.Error("request success with invalid alias_name")
		return
	}
	if err := friendsChatRequest404(handler); err == nil {
		t.Error("request success with invalid path")
		return
	}

	handlerx := FriendsChatPage(ctx, httpLogger, cfg, db, newTsHLSClient(false, true))
	if err := friendsChatRequestPostOK(handlerx); err == nil {
		t.Error("request success with invalid my pubkey")
		return
	}

	handlery := FriendsChatPage(ctx, httpLogger, cfg, newTsDatabase(false, true), newTsHLSClient(true, true))
	if err := friendsChatRequestPostOK(handlery); err == nil {
		t.Error("request success with invalid push message")
		return
	}

	handlerz := FriendsChatPage(ctx, httpLogger, cfg, newTsDatabase(true, false), newTsHLSClient(true, true))
	if err := friendsChatRequestOK(handlerz); err == nil {
		t.Error("request success with invalid load message")
		return
	}
}

func friendsChatRequestHasNotGraphicChars(handler http.HandlerFunc) error {
	formData := url.Values{
		"method":        {"POST"},
		"input_message": {"hello,\x01world!"},
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

func friendsChatRequestInputMessage(handler http.HandlerFunc) error {
	formData := url.Values{
		"method": {"POST"},
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

func friendsChatRequestAliasNameNotFound(handler http.HandlerFunc) error {
	formData := url.Values{
		"method":        {"POST"},
		"input_message": {"hello, world!"},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/friends/chat?alias_name=notfound", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusSeeOther {
		return errors.New("bad status code")
	}

	return nil
}

func friendsChatRequestAliasName(handler http.HandlerFunc) error {
	formData := url.Values{
		"method":        {"POST"},
		"input_message": {"hello, world!"},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/friends/chat?alias_name=", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusSeeOther {
		return errors.New("bad status code")
	}

	return nil
}

func friendsChatRequestPutOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/friends/chat?alias_name=abc", strings.NewReader(`-----------------------------60072964122428079602614660621
Content-Disposition: form-data; name="method"

PUT
-----------------------------60072964122428079602614660621
Content-Disposition: form-data; name="input_file"; filename="file.txt"
Content-Type: text/plain

hello, world!

-----------------------------60072964122428079602614660621--
`))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=---------------------------60072964122428079602614660621")

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusSeeOther {
		return errors.New("bad status code")
	}

	return nil
}

func friendsChatRequestPutSizeFile(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/friends/chat?alias_name=abc", strings.NewReader(`-----------------------------60072964122428079602614660621
Content-Disposition: form-data; name="method"

PUT
-----------------------------60072964122428079602614660621
Content-Disposition: form-data; name="input_file"; filename="file.txt"
Content-Type: text/plain

-----------------------------60072964122428079602614660621--
`))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=---------------------------60072964122428079602614660621")

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusSeeOther {
		return errors.New("bad status code")
	}

	return nil
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
