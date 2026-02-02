package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/number571/go-peer/pkg/logger"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func TestHandleChatMessageAPI(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	httpLogger := std_logger.NewStdLogger(
		func() std_logger.ILogging {
			logging, err := std_logger.LoadLogging([]string{})
			if err != nil {
				panic(err)
			}
			return logging
		}(),
		func(_ logger.ILogArg) string {
			return ""
		},
	)

	handlerX := HandleChatMessageAPI(ctx, httpLogger, &tsConfig{}, newTsHLKClient(true, true, true), newTsDatabase(true, true))
	if err := chatMessageGetRequestOK(handlerX); err != nil {
		t.Fatal(err)
	}

	handlerY := HandleChatMessageAPI(ctx, httpLogger, &tsConfig{}, newTsHLKClient(true, false, true), newTsDatabase(true, true))
	if err := chatMessageGetRequestOK(handlerY); err == nil {
		t.Fatal("success request with get settings error")
	}

	if err := chatMessagePostRequestOK(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := chatMessageRequestInvalidMethod(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := chatMessageRequestInvalidMessage(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := chatMessageRequestWithoutFriend(handlerX); err != nil {
		t.Fatal(err)
	}

	handlerZ := HandleChatMessageAPI(ctx, httpLogger, &tsConfig{}, newTsHLKClient(false, true, true), newTsDatabase(true, true))
	if err := chatMessagePostRequestOK(handlerZ); err == nil {
		t.Fatal("success request with get pub key error")
	}
	handlerA := HandleChatMessageAPI(ctx, httpLogger, &tsConfig{}, newTsHLKClient(true, true, false), newTsDatabase(true, true))
	if err := chatMessagePostRequestOK(handlerA); err == nil {
		t.Fatal("success request with send message error")
	}
	handlerB := HandleChatMessageAPI(ctx, httpLogger, &tsConfig{}, newTsHLKClient(true, true, true), newTsDatabase(true, false))
	if err := chatMessagePostRequestOK(handlerB); err == nil {
		t.Fatal("success request with send message error")
	}
}

func chatMessageGetRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?friend=abc", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func chatMessagePostRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/?friend=abc", bytes.NewBuffer([]byte("hello, world!")))

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func chatMessageRequestInvalidMethod(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusMethodNotAllowed {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func chatMessageRequestInvalidMessage(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("")))

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusBadRequest {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func chatMessageRequestWithoutFriend(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("hello, world!")))

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusForbidden {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}
