// nolint: err113
package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/filesharer/pkg/app/config"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func TestStoragePage(t *testing.T) {
	t.Parallel()

	logging, err := std_logger.LoadLogging([]string{})
	if err != nil {
		t.Error(err)
		return
	}

	httpLogger := std_logger.NewStdLogger(
		logging,
		func(_ logger.ILogArg) string {
			return ""
		},
	)

	cfg := &config.SConfig{
		FSettings: &config.SConfigSettings{
			FLanguage: "ENG",
		},
	}

	ctx := context.Background()
	handler := StoragePage(ctx, httpLogger, cfg, newTsHLSClient(true, true))
	if err := storageRequestOK(handler); err != nil {
		t.Error(err)
		return
	}
	if err := storageRequestDownloadOK(handler); err != nil {
		t.Error(err)
		return
	}

	if err := storageRequest404(handler); err == nil {
		t.Error("request success with invalid path")
		return
	}
	if err := storageRequestAliasName(handler); err == nil {
		t.Error("request success with alias_name")
		return
	}
	handlerx := StoragePage(ctx, httpLogger, cfg, newTsHLSClient(false, true))
	if err := storageRequestOK(handlerx); err == nil {
		t.Error("request success with fetch failed")
		return
	}
}

func storageRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/friends/storage?alias_name=abc", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

func storageRequestAliasName(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/friends/storage?alias_name=&page=0", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

func storageRequestDownloadOK(handler http.HandlerFunc) error {
	fileBytes, err := os.ReadFile("./testdata/file.txt")
	if err != nil {
		return err
	}
	hash := hashing.NewHasher(fileBytes).ToString()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodGet,
		"/friends/storage?alias_name=abc&file_name=file.txt&file_size=13&file_hash="+hash,
		nil,
	)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}

func storageRequest404(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/friends/storage/undefined", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}
