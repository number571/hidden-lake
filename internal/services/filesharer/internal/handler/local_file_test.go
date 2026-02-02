package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/number571/go-peer/pkg/logger"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func TestHandleLocalFileAPI(t *testing.T) {
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

	handlerX := HandleLocalFileAPI(ctx, httpLogger, newTsHLKClient(0, true), "./testdata")
	if err := localFileRequestGetOK(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := localFileRequestPostOK(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := localFileRequestDeleteOK(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := localFileRequestInvalidMethod(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := localFileRequestFileNotFound(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := localFileRequestInvalidFileName(handlerX); err != nil {
		t.Fatal(err)
	}
	if err := localFileRequestPersonalNotFound(handlerX); err != nil {
		t.Fatal(err)
	}
}

func localFileRequestGetOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?friend=abc&name=something.txt", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if !bytes.Equal(body, []byte(`something111
222
333
`)) {
		return errors.New("error") // nolint: err113
	}

	return nil
}

func localFileRequestPersonalNotFound(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?name=file.txt&friend=111", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusForbidden {
		fmt.Println(res.StatusCode)
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func localFileRequestInvalidFileName(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?name=", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusBadRequest {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func localFileRequestPostOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/?friend=abc&name=something111.txt", bytes.NewBuffer([]byte("571_175")))

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func localFileRequestDeleteOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/?friend=abc&name=something111.txt", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func localFileRequestFileNotFound(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?friend=abc&name=something111.txt", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusNotFound {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func localFileRequestInvalidMethod(handler http.HandlerFunc) error {
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
