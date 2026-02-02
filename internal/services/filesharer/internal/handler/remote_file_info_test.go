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

func TestHandleRemoteFileInfoAPI(t *testing.T) {
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

	handlerX := HandleRemoteFileInfoAPI(ctx, httpLogger, newTsHLKClient(0, true))
	if err := remoteFileInfoRequestOK(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := remoteFileInfoRequestInvalidMethod(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := remoteFileInfoRequestInvalidPersonal(handlerX); err != nil {
		t.Fatal(err)
	}

	handlerY := HandleRemoteFileInfoAPI(ctx, httpLogger, newTsHLKClient(1, true))
	if err := remoteFileInfoRequestOK(handlerY); err == nil {
		t.Fatal("success request with fetch error")
	}

	handlerZ := HandleRemoteFileInfoAPI(ctx, httpLogger, newTsHLKClient(-1, true))
	if err := remoteFileInfoRequestOK(handlerZ); err == nil {
		t.Fatal("success request with response status error")
	}

	handlerA := HandleRemoteFileInfoAPI(ctx, httpLogger, newTsHLKClient(-2, true))
	if err := remoteFileInfoRequestOK(handlerA); err == nil {
		t.Fatal("success request with invalid response")
	}

	handlerB := HandleRemoteFileInfoAPI(ctx, httpLogger, newTsHLKClient(-3, true))
	if err := remoteFileInfoRequestOK(handlerB); err == nil {
		t.Fatal("success request with invalid file name")
	}
}

func remoteFileInfoRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?friend=abc&name=example.txt&personal", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func remoteFileInfoRequestInvalidMethod(handler http.HandlerFunc) error {
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

func remoteFileInfoRequestInvalidPersonal(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?friend=abc&name=example.txt&personal=8492", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusBadRequest {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}
