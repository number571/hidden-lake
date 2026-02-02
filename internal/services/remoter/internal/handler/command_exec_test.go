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

func TestHandleCommandExecAPI(t *testing.T) {
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

	handlerX := HandleCommandExecAPI(ctx, httpLogger, newTsHLKClient(0))
	if err := commandExecRequestOK(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := commandExecRequestInvalidMethod(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := commandExecRequestInvalidBody(handlerX); err != nil {
		t.Fatal(err)
	}

	handlerY := HandleCommandExecAPI(ctx, httpLogger, newTsHLKClient(-1))
	if err := commandExecRequestOK(handlerY); err == nil {
		t.Fatal("success with invalid response (1)")
	}

	handlerA := HandleCommandExecAPI(ctx, httpLogger, newTsHLKClient(1))
	if err := commandExecRequestOK(handlerA); err == nil {
		t.Fatal("success with invalid fetch")
	}
}

func commandExecRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(`{"password":"","command":[]}`)))

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func commandExecRequestInvalidMethod(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusMethodNotAllowed {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func commandExecRequestInvalidBody(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(`AAAbbb111`)))

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusConflict {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}
