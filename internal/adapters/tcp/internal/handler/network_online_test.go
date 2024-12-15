package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/number571/go-peer/pkg/logger"
)

func TestHandleNetworkOnlineAPI(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	handler := HandleNetworkOnlineAPI(log, &tsNetworkNode{})
	if err := networkOnlineRequestMethod(handler); err != nil {
		t.Error(err)
		return
	}
	if err := networkOnlineRequestGET(handler); err != nil {
		t.Error(err)
		return
	}
	if err := networkOnlineRequestURLParse(handler); err != nil {
		t.Error(err)
		return
	}
	if err := networkOnlineRequestURLScheme(handler); err != nil {
		t.Error(err)
		return
	}
	if err := networkOnlineRequestDelConnection(handler, http.StatusOK); err != nil {
		t.Error(err)
		return
	}

	handlerx := HandleNetworkOnlineAPI(log, &tsNetworkNode{fWithFail: true})
	if err := networkOnlineRequestDelConnection(handlerx, http.StatusInternalServerError); err != nil {
		t.Error(err)
		return
	}
}

func networkOnlineRequestMethod(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

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
	defer res.Body.Close()

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
	defer res.Body.Close()

	if res.StatusCode != http.StatusTeapot {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func networkOnlineRequestURLScheme(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", strings.NewReader("http://abc"))

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func networkOnlineRequestDelConnection(handler http.HandlerFunc, code int) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", strings.NewReader("tcp://abc"))

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != code {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}
