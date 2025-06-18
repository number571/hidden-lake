package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/pinger/pkg/app/config"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func TestHandleIncomingPingHTTP(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)
	handler := HandleIncomingPingHTTP(&tsConfig{}, log)

	if err := incomingRequestMethod(handler); err != nil {
		t.Fatal(err)
	}
	if err := incomingRequestSuccess(handler); err != nil {
		t.Fatal(err)
	}
}

func incomingRequestMethod(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusMethodNotAllowed {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingRequestSuccess(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

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

func (p *tsConfigSettings) GetResponseMessage() string { return "" }
