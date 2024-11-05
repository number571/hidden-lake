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

	handler := HandleConfigSettingsAPI(newTsWrapper(true), httpLogger, newTsNode(true, true, true, true))
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

	addr := testutils.TgAddrs[16]
	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 2)
	pathDB := fmt.Sprintf(tcPathDBTemplate, 2)

	_, node, _, cancel, srv := testAllCreate(pathCfg, pathDB, addr)
	defer testAllFree(node, cancel, srv, pathCfg, pathDB)

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			"http://"+addr,
			&http.Client{Timeout: time.Minute},
		),
	)

	sett, err := client.GetSettings(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if sett.GetQueuePeriodMS() != 1000 {
		t.Error("invalid queue period")
		return
	}

	if sett.GetMessageSizeBytes() != (8 << 10) {
		t.Error("invalid message size")
		return
	}

	if sett.GetWorkSizeBits() != 22 {
		t.Error("invalid work size")
		return
	}
}
