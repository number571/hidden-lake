package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/logger"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestHandleConfigSettingsAPI2(t *testing.T) {
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

	handler := HandleConfigSettingsAPI(&tsConfig{}, httpLogger)
	if err := settingsAPIRequestOK(handler); err != nil {
		t.Error(err)
		return
	}

	if err := settingsAPIRequestMethod(handler); err == nil {
		t.Error("request success with invalid method")
		return
	}
}

func settingsAPIRequestOK(handler http.HandlerFunc) error {
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

func settingsAPIRequestMethod(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)

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

func TestHandleConfigSettingsAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[22]
	os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))

	srv, cancel, db, hltClient := testAllRun(addr)
	defer testAllFree(addr, srv, cancel, db)

	settings, err := hltClient.GetSettings(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if settings.GetNetworkKey() != tcNetworkKey {
		t.Error("incorrect network key")
		return
	}

	if settings.GetWorkSizeBits() != tcWorkSize {
		t.Error("incorrect work size bits")
		return
	}

	if settings.GetMessagesCapacity() != tcCapacity {
		t.Error("incorrect messages capacity")
		return
	}

	if settings.GetMessageSizeBytes() != tcMessageSize {
		t.Error("incorrect messages size bytes")
		return
	}
}
