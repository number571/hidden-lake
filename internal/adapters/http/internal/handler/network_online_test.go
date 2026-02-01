package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/message/layer1"
	hla_http "github.com/number571/hidden-lake/pkg/network/adapters/http"
)

func TestHandleNetworkOnlineAPI(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	handler := HandleNetworkOnlineAPI(log, &tsHTTPAdapter{})
	if err := networkOnlineRequestMethod(handler); err != nil {
		t.Fatal(err)
	}
	if err := networkOnlineRequestGET(handler); err != nil {
		t.Fatal(err)
	}
	if err := networkOnlineRequestURLParse(handler); err != nil {
		t.Fatal(err)
	}
	if err := networkOnlineRequestURLScheme(handler); err != nil {
		t.Fatal(err)
	}
	if err := networkOnlineRequestDelConnection(handler, http.StatusNoContent); err != nil {
		t.Fatal(err)
	}
}

func networkOnlineRequestMethod(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusMethodNotAllowed {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func networkOnlineRequestGET(handler http.HandlerFunc) error {
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

func networkOnlineRequestURLParse(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", strings.NewReader("94.103.91.81:9581"))

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusTeapot {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func networkOnlineRequestURLScheme(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", strings.NewReader("tcp://abc"))

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusAccepted {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func networkOnlineRequestDelConnection(handler http.HandlerFunc, code int) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", strings.NewReader("http://abc"))

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != code {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

type tsHTTPAdapter struct{}

func (p *tsHTTPAdapter) Run(context.Context) error                        { return nil }
func (p *tsHTTPAdapter) Produce(context.Context, layer1.IMessage) error   { return nil }
func (p *tsHTTPAdapter) Consume(context.Context) (layer1.IMessage, error) { return nil, nil }

func (p *tsHTTPAdapter) WithLogger(string, logger.ILogger) hla_http.IHTTPAdapter {
	return p
}
func (p *tsHTTPAdapter) WithHandlers(map[string]http.HandlerFunc) hla_http.IHTTPAdapter { return p }
func (p *tsHTTPAdapter) GetOnlines() []string {
	return []string{"127.0.0.1", "127.0.0.2"}
}
