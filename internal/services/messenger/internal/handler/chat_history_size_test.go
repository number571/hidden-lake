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

func TestHandleChatHistorySizeAPI(t *testing.T) {
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

	handlerX := HandleChatHistorySizeAPI(ctx, httpLogger, &tsConfig{}, newTsHLKClient(true, true, true), newTsDatabase(true, true))
	if err := chatHistorySizeRequestOK(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := chatHistorySizeRequestInvalidMethod(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := chatHistorySizeRequestWithoutFriend(handlerX); err != nil {
		t.Fatal(err)
	}

	handlerY := HandleChatHistorySizeAPI(ctx, httpLogger, &tsConfig{}, newTsHLKClient(false, true, true), newTsDatabase(true, true))
	if err := chatHistorySizeRequestOK(handlerY); err == nil {
		t.Fatal("success request with get pub key error")
	}
}

func chatHistorySizeRequestOK(handler http.HandlerFunc) error {
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

func chatHistorySizeRequestInvalidMethod(handler http.HandlerFunc) error {
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

func chatHistorySizeRequestWithoutFriend(handler http.HandlerFunc) error {
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
