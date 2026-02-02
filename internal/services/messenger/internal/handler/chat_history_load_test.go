package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/number571/go-peer/pkg/logger"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func TestHandleChatHistoryLoadAPI(t *testing.T) {
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

	handlerX := HandleChatHistoryLoadAPI(ctx, httpLogger, &tsConfig{}, newTsHLKClient(true, true, true), newTsDatabase(true, true))
	if err := chatHistoryLoadRequestOK(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := chatHistoryLoadRequestInvalidMethod(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := chatHistoryLoadRequestWithoutFriend(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := chatHistoryLoadRequestCountParse(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := chatHistoryLoadRequestStartParse(handlerX); err != nil {
		t.Fatal(err)
	}

	handlerY := HandleChatHistoryLoadAPI(ctx, httpLogger, &tsConfig{}, newTsHLKClient(false, true, true), newTsDatabase(true, true))
	if err := chatHistoryLoadRequestOK(handlerY); err == nil {
		t.Fatal("success request with get pub key error")
	}

	handlerZ := HandleChatHistoryLoadAPI(ctx, httpLogger, &tsConfig{}, newTsHLKClient(true, true, true), newTsDatabase(false, true))
	if err := chatHistoryLoadRequestOK(handlerZ); err == nil {
		t.Fatal("success request with load message error")
	}
}

func chatHistoryLoadRequestOK(handler http.HandlerFunc) error {
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

func chatHistoryLoadRequestInvalidMethod(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusMethodNotAllowed {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func chatHistoryLoadRequestWithoutFriend(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusForbidden {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func chatHistoryLoadRequestCountParse(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?friend=abc&count=qqq", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusBadRequest {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func chatHistoryLoadRequestStartParse(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?friend=abc&start=qqq", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusBadRequest {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}
