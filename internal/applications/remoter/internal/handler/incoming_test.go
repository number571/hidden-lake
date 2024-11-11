package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/remoter/pkg/app/config"
	"github.com/number571/hidden-lake/internal/applications/remoter/pkg/client"
	"github.com/number571/hidden-lake/internal/applications/remoter/pkg/settings"
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestIncomingExecHTTP2(t *testing.T) {
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

	ctx := context.Background()
	cfg := &config.SConfig{
		FSettings: &config.SConfigSettings{
			FExecTimeoutMS: 10_000,
			FPassword:      tcPassword,
		},
	}

	handler := HandleIncomingExecHTTP(ctx, cfg, httpLogger)
	if err := incomingExecRequestOK(handler); err != nil {
		t.Error(err)
		return
	}

	if err := incomingExecRequestMethod(handler); err == nil {
		t.Error("success request with invalid method")
		return
	}
	if err := incomingExecRequestPassword(handler); err == nil {
		t.Error("success request with invalid password")
		return
	}
	if err := incomingExecRequestHasNotGraphicChars(handler); err == nil {
		t.Error("success request with invalid chars (not graphic)")
		return
	}
	if err := incomingExecRequestCommand(handler); err == nil {
		t.Error("success request with invalid command")
		return
	}
}

func incomingExecRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/exec", strings.NewReader("echo"+settings.CExecSeparator+"hello"))
	req.Header.Set(settings.CHeaderPassword, tcPassword)

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

func incomingExecRequestCommand(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/exec", strings.NewReader("____"+settings.CExecSeparator+"hello"))
	req.Header.Set(settings.CHeaderPassword, tcPassword)

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

func incomingExecRequestHasNotGraphicChars(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/exec", strings.NewReader("echo"+settings.CExecSeparator+"\x01hello"))
	req.Header.Set(settings.CHeaderPassword, tcPassword)

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

func incomingExecRequestPassword(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/exec", strings.NewReader("echo"+settings.CExecSeparator+"hello"))
	req.Header.Set(settings.CHeaderPassword, tcPassword+"_")

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

func incomingExecRequestMethod(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/exec", strings.NewReader("echo"+settings.CExecSeparator+"hello"))
	req.Header.Set(settings.CHeaderPassword, tcPassword)

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

func TestIncomingExecHTTP1(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, service := testRunService(testutils.TgAddrs[40])
	defer service.Close()

	testRunNewNodes(
		ctx,
		testutils.TgAddrs[41],
		testutils.TgAddrs[42],
		testutils.TgAddrs[40],
	)

	hlsClient := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			"http://"+testutils.TgAddrs[41],
			&http.Client{Timeout: time.Minute},
		),
	)

	hlrClient := client.NewClient(
		client.NewBuilder(tcPassword),
		client.NewRequester(hlsClient),
	)

	msg := "hello, world!"
	rsp, err := hlrClient.Exec(ctx, "test_recv", "echo", msg)
	if err != nil {
		t.Error(err)
		return
	}

	if strings.TrimSpace(string(rsp)) != msg {
		t.Error("get invalid response")
		return
	}
}
