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
	"time"

	"github.com/number571/go-peer/pkg/logger"
	hle_client "github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/client"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestHandlePubKeyAPI2(t *testing.T) {
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

	handler := HandleServicePubKeyAPI(httpLogger, tgPrivKey1.GetPubKey())
	if err := pubkeyAPIRequestOK(handler); err != nil {
		t.Error(err)
		return
	}

	if err := pubkeyAPIRequestMethod(handler); err == nil {
		t.Error("request success with invalid method")
		return
	}
}

func pubkeyAPIRequestOK(handler http.HandlerFunc) error {
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

func pubkeyAPIRequestMethod(handler http.HandlerFunc) error {
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

func TestHandlePubKeyAPI(t *testing.T) {
	t.Parallel()

	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 3)
	defer os.Remove(pathCfg)

	_, service := testRunService(pathCfg, testutils.TgAddrs[31])
	defer service.Close()

	time.Sleep(100 * time.Millisecond)
	hleClient := hle_client.NewClient(
		hle_client.NewBuilder(),
		hle_client.NewRequester(
			testutils.TgAddrs[31],
			&http.Client{Timeout: time.Second / 2},
			testNetworkMessageSettings(),
		),
	)

	gotPubKey, err := hleClient.GetPubKey(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	pubKey := tgPrivKey1.GetPubKey()
	if pubKey.ToString() != gotPubKey.ToString() {
		t.Error("public keys not equals")
		return
	}
}
