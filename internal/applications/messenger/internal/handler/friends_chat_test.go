// nolint: err113
package handler

import (
	"bytes"
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
		t.Fatal("success get message limit with settings failed")
	}

	hlsClient := newTsHLSClient(true, false)
	hlsClient.fPldSize = 10
	if _, err := utils.GetMessageLimit(ctx, hlsClient); err == nil {
		t.Fatal("success get message limit with payload limit")
	}
}

func wrapText(pMsg string) []byte {
	return bytes.Join([][]byte{
		{0x01}, // cIsText
		[]byte(pMsg),
	}, []byte{})
}

func TestFriendsChatPage(t *testing.T) {
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

	db := newTsDatabase(true, true)
	_ = db.Push(nil, database.NewMessage(false, wrapText("hello, world!")))

	handler := FriendsChatPage(ctx, httpLogger, cfg, db, newTsHLSClient(true, true))
	if err := friendsChatRequestOK(handler); err != nil {
		t.Fatal(err)
	}
	if err := friendsChatRequestPostOK(handler); err != nil {
		t.Fatal(err)
	}
	if err := friendsChatRequestPutOK(handler); err != nil {
		t.Fatal(err)
	}
	if err := friendsChatRequestPostPingOK(handler); err != nil {
		t.Fatal(err)
	}

	if err := friendsChatRequestPutSizeFile(handler); err == nil {
		t.Fatal("request success with invalid file size")
	}
	if err := friendsChatRequestHasNotGraphicChars(handler); err == nil {
		t.Fatal("request success with invalid (has not grpahic chars)")
	}
	if err := friendsChatRequestInputMessage(handler); err == nil {
		t.Fatal("request success with invalid input_message")
	}
	if err := friendsChatRequestAliasNameNotFound(handler); err == nil {
		t.Fatal("request success with not found alias_name")
	}
	if err := friendsChatRequestAliasName(handler); err == nil {
		t.Fatal("request success with invalid alias_name")
	}
	if err := friendsChatRequest404(handler); err == nil {
		t.Fatal("request success with invalid path")
	}

	handlerx := FriendsChatPage(ctx, httpLogger, cfg, db, newTsHLSClient(false, true))
	if err := friendsChatRequestPostOK(handlerx); err == nil {
		t.Fatal("request success with invalid my pubkey")
	}

	handlery := FriendsChatPage(ctx, httpLogger, cfg, newTsDatabase(false, true), newTsHLSClient(true, true))
	if err := friendsChatRequestPostOK(handlery); err == nil {
		t.Fatal("request success with invalid push message")
	}

	handlerz := FriendsChatPage(ctx, httpLogger, cfg, newTsDatabase(true, false), newTsHLSClient(true, true))
	if err := friendsChatRequestOK(handlerz); err == nil {
		t.Fatal("request success with invalid load message")
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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}
