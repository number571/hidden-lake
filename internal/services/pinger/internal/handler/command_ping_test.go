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

func TestHandleCommandPingAPI(t *testing.T) {
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

	handlerX := HandleCommandPingAPI(ctx, httpLogger, newTsHLKClient(0))
	if err := commandPingRequestOK(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := commandPingRequestInvalidMethod(handlerX); err != nil {
		t.Fatal(err)
	}

	handlerY := HandleCommandPingAPI(ctx, httpLogger, newTsHLKClient(-1))
	if err := commandPingRequestOK(handlerY); err == nil {
		t.Fatal("success with invalid response (1)")
	}

	handlerZ := HandleCommandPingAPI(ctx, httpLogger, newTsHLKClient(-2))
	if err := commandPingRequestOK(handlerZ); err == nil {
		t.Fatal("success with invalid response (2)")
	}

	handlerA := HandleCommandPingAPI(ctx, httpLogger, newTsHLKClient(1))
	if err := commandPingRequestOK(handlerA); err == nil {
		t.Fatal("success with invalid fetch")
	}
}

func commandPingRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func commandPingRequestInvalidMethod(handler http.HandlerFunc) error {
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
