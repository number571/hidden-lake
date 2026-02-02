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

func TestHandleRemoteFileAPI(t *testing.T) {
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

	handlerX := HandleRemoteFileAPI(ctx, &tsConfig{}, httpLogger, newTsHLKClient(0, true), "./testdata")
	if err := remoteFileRequestOK(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := remoteFileRequestInvalidMethod(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := remoteFileRequestInvalidPersonal(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := remoteFileRequestNotFoundFriend(handlerX); err != nil {
		t.Fatal(err)
	}

	handlerY := HandleRemoteFileAPI(ctx, &tsConfig{}, httpLogger, newTsHLKClient(1, true), "./testdata")
	if err := remoteFileRequestOK(handlerY); err == nil {
		t.Fatal("success request with invalid fetch")
	}

	handlerZ := HandleRemoteFileAPI(ctx, &tsConfig{}, httpLogger, newTsHLKClient(-1, true), "./testdata")
	if err := remoteFileRequestOK(handlerZ); err == nil {
		t.Fatal("success request with invalid status code")
	}

	handlerA := HandleRemoteFileAPI(ctx, &tsConfig{}, httpLogger, newTsHLKClient(-2, true), "./testdata")
	if err := remoteFileRequestOK(handlerA); err == nil {
		t.Fatal("success request with invalid response body")
	}

	handlerB := HandleRemoteFileAPI(ctx, &tsConfig{}, httpLogger, newTsHLKClient(-3, true), "./testdata")
	if err := remoteFileRequestOK(handlerB); err == nil {
		t.Fatal("success request with invalid file name")
	}

	handlerC := HandleRemoteFileAPI(ctx, &tsConfig{}, httpLogger, newTsHLKClient(0, false), "./testdata")
	if err := remoteFileRequestOK(handlerC); err == nil {
		t.Fatal("success request with build stream error")
	}
}

func remoteFileRequestOK(handler http.HandlerFunc) error {
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

func remoteFileRequestNotFoundFriend(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?friend=qwerty&name=example.txt&personal", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusForbidden {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func remoteFileRequestInvalidPersonal(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?friend=abc&name=example.txt&personal=7492", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusBadRequest {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func remoteFileRequestInvalidMethod(handler http.HandlerFunc) error {
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
