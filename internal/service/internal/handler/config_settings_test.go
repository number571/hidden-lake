package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
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

	handler := HandleConfigSettingsAPI(newTsWrapper(true), httpLogger, newTsNode(true, true, true))
	if err := settingsAPIRequestOK(handler); err != nil {
		t.Fatal(err)
	}

	if err := settingsAPIRequestMethod(handler); err == nil {
		t.Fatal("request success with invalid method")
	}
}

func settingsAPIRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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

	addr := testutils.TgAddrs[16]
	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 2)
	pathDB := fmt.Sprintf(tcPathDBTemplate, 2)

	_, node, _, cancel, srv := testAllCreate(pathCfg, pathDB, addr)
	defer testAllFree(node, cancel, srv, pathCfg, pathDB)

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			addr,
			&http.Client{Timeout: time.Minute},
		),
	)

	sett, err := client.GetSettings(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if sett.GetQueuePeriod() != time.Second {
		t.Fatal("invalid queue period")
	}

	if sett.GetMessageSizeBytes() != (8 << 10) {
		t.Fatal("invalid message size")
	}

	if sett.GetWorkSizeBits() != 22 {
		t.Fatal("invalid work size")
	}
}
