// nolint: goerr113
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
	"github.com/number571/hidden-lake/internal/applications/notifier/internal/database"
	"github.com/number571/hidden-lake/internal/applications/notifier/internal/utils"
	"github.com/number571/hidden-lake/internal/applications/notifier/pkg/app/config"
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

func wrapText(pMsg string) []byte {
	return bytes.Join([][]byte{
		{0x01}, // cIsText
		[]byte(pMsg),
	}, []byte{})
}

func TestChannelsChatPage(t *testing.T) {
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

	handler := ChannelsChatPage(ctx, httpLogger, cfg, db, newTsHLSClient(true, true))
	if err := channelsChatRequestOK(handler); err != nil {
		t.Error(err)
		return
	}
	if err := channelsChatRequestPostOK(handler); err != nil {
		t.Error(err)
		return
	}
	if err := channelsChatRequestPutOK(handler); err != nil {
		t.Error(err)
		return
	}

	if err := channelsChatRequestPutSizeFile(handler); err == nil {
		t.Error("request success with invalid file size")
		return
	}
	if err := channelsChatRequestHasNotGraphicChars(handler); err == nil {
		t.Error("request success with invalid (has not grpahic chars)")
		return
	}
	if err := channelsChatRequestInputMessage(handler); err == nil {
		t.Error("request success with invalid input_message")
		return
	}
	if err := channelsChatRequestAliasName(handler); err == nil {
		t.Error("request success with invalid alias_name")
		return
	}
	if err := channelsChatRequest404(handler); err == nil {
		t.Error("request success with invalid path")
		return
	}

	handlerx := ChannelsChatPage(ctx, httpLogger, cfg, db, newTsHLSClient(false, true))
	if err := channelsChatRequestPostOK(handlerx); err == nil {
		t.Error("request success with invalid my pubkey")
		return
	}

	handlery := ChannelsChatPage(ctx, httpLogger, cfg, newTsDatabase(false, true), newTsHLSClient(true, true))
	if err := channelsChatRequestPostOK(handlery); err == nil {
		t.Error("request success with invalid push message")
		return
	}

	handlerz := ChannelsChatPage(ctx, httpLogger, cfg, newTsDatabase(true, false), newTsHLSClient(true, true))
	if err := channelsChatRequestOK(handlerz); err == nil {
		t.Error("request success with invalid load message")
		return
	}
}

func channelsChatRequestHasNotGraphicChars(handler http.HandlerFunc) error {
	formData := url.Values{
		"method":        {"POST"},
		"input_message": {"hello,\x01world!"},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/channels/chat?key=abc", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusSeeOther {
		return errors.New("bad status code")
	}

	return nil
}

func channelsChatRequestInputMessage(handler http.HandlerFunc) error {
	formData := url.Values{
		"method": {"POST"},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/channels/chat?key=abc", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusSeeOther {
		return errors.New("bad status code")
	}

	return nil
}

func channelsChatRequestAliasName(handler http.HandlerFunc) error {
	formData := url.Values{
		"method":        {"POST"},
		"input_message": {"hello, world!"},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/channels/chat?key=", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusSeeOther {
		return errors.New("bad status code")
	}

	return nil
}

func channelsChatRequestPutOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/channels/chat?key=abc", strings.NewReader(`-----------------------------60072964122428079602614660621
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

func channelsChatRequestPutSizeFile(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/channels/chat?key=abc", strings.NewReader(`-----------------------------60072964122428079602614660621
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

func channelsChatRequestPostOK(handler http.HandlerFunc) error {
	formData := url.Values{
		"method":        {"POST"},
		"input_message": {"hello, world!"},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/channels/chat?key=abc", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusSeeOther {
		return errors.New("bad status code")
	}

	return nil
}

func channelsChatRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/channels/chat?key=abc", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

func channelsChatRequest404(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/channels/chat/undefined", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}
