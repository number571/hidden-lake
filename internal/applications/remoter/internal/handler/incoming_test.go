package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/remoter/pkg/app/config"
	hlr_settings "github.com/number571/hidden-lake/internal/applications/remoter/pkg/settings"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func TestHandleIncomingExecHTTP(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	log := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)
	handler := HandleIncomingExecHTTP(ctx, &tsConfig{}, log)

	if err := incomingRequestMethod(handler); err != nil {
		t.Fatal(err)
	}
	if err := incomingRequestInvalidPassword(handler); err != nil {
		t.Fatal(err)
	}
	if err := incomingRequestNotGraphicChars(handler); err != nil {
		t.Fatal(err)
	}
	if err := incomingRequestExecCommand(handler); err != nil {
		t.Fatal(err)
	}
	if err := incomingRequestSuccess(handler); err != nil {
		t.Fatal(err)
	}
}

func incomingRequestMethod(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusMethodNotAllowed {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingRequestInvalidPassword(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(hlr_settings.CHeaderPassword, "___")

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusForbidden {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingRequestNotGraphicChars(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("abc\x01xyz")))
	req.Header.Set(hlr_settings.CHeaderPassword, "abc")

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusBadRequest {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingRequestExecCommand(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("")))
	req.Header.Set(hlr_settings.CHeaderPassword, "abc")

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusInternalServerError {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingRequestSuccess(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("dir")))
	req.Header.Set(hlr_settings.CHeaderPassword, "abc")

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

type tsConfig struct{}

func (p *tsConfig) GetSettings() config.IConfigSettings { return &tsConfigSettings{} }
func (p *tsConfig) GetAddress() config.IAddress         { return nil }
func (p *tsConfig) GetLogging() std_logger.ILogging     { return nil }

type tsConfigSettings struct{}

func (p *tsConfigSettings) GetExecTimeout() time.Duration { return time.Second }
func (p *tsConfigSettings) GetPassword() string           { return "abc" }
