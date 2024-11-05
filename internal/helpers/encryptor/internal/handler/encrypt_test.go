package handler

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	hle_settings "github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/settings"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func TestHandleMessageEncryptAPI(t *testing.T) {
	t.Parallel()

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

	handler1 := HandleMessageEncryptAPI(&tsConfig{}, httpLogger, client.NewClient(tgPrivKey1, 8192), 1)
	if err := encryptAPIRequestOK(handler1); err != nil {
		t.Error(err)
		return
	}

	handler := HandleMessageEncryptAPI(&tsConfig{}, httpLogger, newTsClient(true), 1)
	if err := encryptAPIRequestOK(handler); err != nil {
		t.Error(err)
		return
	}

	if err := encryptAPIRequestHexDecode(handler); err == nil {
		t.Error("request success with invalid hex decode")
		return
	}
	if err := encryptAPIRequestNotFound(handler); err == nil {
		t.Error("request success with not found alias_name")
		return
	}
	if err := encryptAPIRequestDecode(handler); err == nil {
		t.Error("request success with invalid decode")
		return
	}
	if err := encryptAPIRequestMethod(handler); err == nil {
		t.Error("request success with invalid method")
		return
	}

	handlerx := HandleMessageEncryptAPI(&tsConfig{}, httpLogger, newTsClient(false), 1)
	if err := encryptAPIRequestOK(handlerx); err == nil {
		t.Error("request success with encrypt failed")
		return
	}
}

func encryptAPIRequestOK(handler http.HandlerFunc) error {
	container := hle_settings.SContainer{
		FAliasName: "abc",
		FPldHead:   1,
		FHexData:   encoding.HexEncode([]byte("hello, world!")),
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(encoding.SerializeJSON(container)))

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func encryptAPIRequestHexDecode(handler http.HandlerFunc) error {
	container := hle_settings.SContainer{
		FAliasName: "abc",
		FPldHead:   1,
		FHexData:   "hello, world!",
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(encoding.SerializeJSON(container)))

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func encryptAPIRequestNotFound(handler http.HandlerFunc) error {
	container := hle_settings.SContainer{
		FAliasName: "notfound",
		FPldHead:   1,
		FHexData:   encoding.HexEncode([]byte("hello, world!")),
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(encoding.SerializeJSON(container)))

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func encryptAPIRequestDecode(handler http.HandlerFunc) error {
	container := hle_settings.SContainer{
		FAliasName: "abc",
		FPldHead:   1,
		FHexData:   encoding.HexEncode([]byte("hello, world!")),
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bytes.Join(
		[][]byte{
			[]byte{1},
			encoding.SerializeJSON(container),
		},
		[]byte{},
	)))

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func encryptAPIRequestMethod(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}
