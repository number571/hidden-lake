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

func TestHandleRemoteListAPI(t *testing.T) {
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

	handlerX := HandleRemoteListAPI(ctx, httpLogger, newTsHLKClient(2, true))
	if err := remoteListRequestOK(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := remoteListRequestInvalidMethod(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := remoteListRequestInvalidPage(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := remoteListRequestInvalidPersonal(handlerX); err != nil {
		t.Fatal(err)
	}

	handlerY := HandleRemoteListAPI(ctx, httpLogger, newTsHLKClient(1, true))
	if err := remoteListRequestOK(handlerY); err == nil {
		t.Fatal("success request with fetch response error")
	}

	handlerZ := HandleRemoteListAPI(ctx, httpLogger, newTsHLKClient(-1, true))
	if err := remoteListRequestOK(handlerZ); err == nil {
		t.Fatal("success request with invalid status code")
	}

	handlerA := HandleRemoteListAPI(ctx, httpLogger, newTsHLKClient(-2, true))
	if err := remoteListRequestOK(handlerA); err == nil {
		t.Fatal("success request with invalid response body")
	}
}

func remoteListRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?friend=abc&personal", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func remoteListRequestInvalidPage(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?friend=abc&personal&page=qqqq", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusBadRequest {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func remoteListRequestInvalidPersonal(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?page=0&personal=gjrid", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusBadRequest {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func remoteListRequestInvalidMethod(handler http.HandlerFunc) error {
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
