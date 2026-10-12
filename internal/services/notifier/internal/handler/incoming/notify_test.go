package incoming

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/number571/go-peer/pkg/logger"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/broker"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func TestHandleIncomingNotifyHTTP(t *testing.T) {
	t.Parallel()

	logging, err := std_logger.LoadLogging([]string{})
	if err != nil {
		t.Fatal(err)
	}

	httpLogger := std_logger.NewStdLogger(
		logging,
		func(_ logger.ILogArg) string {
			return ""
		},
	)

	ctx := context.Background()
	msgBroker := broker.NewDataBroker(256, 256)
	handler := HandleIncomingNotifyHTTP(ctx, httpLogger, msgBroker)

	if err := incomingNotifyRequestOK(handler); err != nil {
		t.Fatal(err)
	}

	if err := incomingNotifyRequestMethod(handler); err == nil {
		t.Fatal("request success with invalid method")
	}
	if err := incomingNotifyRequestPubKey(handler); err == nil {
		t.Fatal("request success with invalid pubkey")
	}
	if err := incomingNotifyRequestMessage(handler); err == nil {
		t.Fatal("request success with invalid message")
	}
}

func incomingNotifyRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewBuffer([]byte("hello, world!")))
	req.Header.Set(hlk_settings.CHeaderSenderName, "abc")

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

func incomingNotifyRequestMessage(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/notify", nil)
	req.Header.Set(hlk_settings.CHeaderSenderName, "abc")

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

func incomingNotifyRequestPubKey(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewBuffer([]byte("hello, world!")))

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

func incomingNotifyRequestMethod(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/notify", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}
