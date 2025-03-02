package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/hidden-lake/internal/applications/notifier/pkg/app/config"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/layer3"
)

func TestHandleIncomingRedirectHTTP(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	log := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	hlsClient := newTsHLSClient(true, true)
	handler := HandleIncomingRedirectHTTP(
		ctx,
		&config.SConfig{FSettings: &config.SConfigSettings{}},
		log,
		hlsClient,
	)

	if err := incomingRedirectRequestMethod(handler); err != nil {
		t.Error(err)
		return
	}
	if err := incomingRedirectRequestLoadMessage(handler); err != nil {
		t.Error(err)
		return
	}
	if err := incomingRedirectRequestExtractMessage(handler); err != nil {
		t.Error(err)
		return
	}

	hlsClient2 := newTsHLSClient(true, true)
	hlsClient2.fFriends = map[string]asymmetric.IPubKey{}
	handler2 := HandleIncomingRedirectHTTP(
		ctx,
		&config.SConfig{FSettings: &config.SConfigSettings{}},
		log,
		hlsClient2,
	)
	if err := incomingRedirectRequestGetFriends(handler2); err != nil {
		t.Error(err)
		return
	}

	if err := incomingRedirectRequestLoadPubKey(handler); err != nil {
		t.Error(err)
		return
	}
	if err := incomingRedirectRequestGetAliasByPubKey(handler); err != nil {
		t.Error(err)
		return
	}

	hlsClient3 := newTsHLSClient(true, false)
	handler3 := HandleIncomingRedirectHTTP(
		ctx,
		&config.SConfig{FSettings: &config.SConfigSettings{}},
		log,
		hlsClient3,
	)
	if err := incomingRedirectRequestRedirect(handler3, hlsClient3.fFriendPubKey); err != nil {
		t.Error(err)
		return
	}

	if err := incomingRedirectRequestOK(handler, hlsClient.fFriendPubKey); err != nil {
		t.Error(err)
		return
	}
}

func incomingRedirectRequestMethod(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusMethodNotAllowed {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingRedirectRequestLoadMessage(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte{1, 2, 3}))

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingRedirectRequestExtractMessage(handler http.HandlerFunc) error {
	msg := layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{}),
		}),
		payload.NewPayload32(1, []byte("hello")),
	)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(msg.ToBytes()))

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotAcceptable {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingRedirectRequestGetFriends(handler http.HandlerFunc) error {
	msg := layer3.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{}),
		}),
		[]byte("hello, world!"),
	)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(msg.ToBytes()))

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingRedirectRequestLoadPubKey(handler http.HandlerFunc) error {
	msg := layer3.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{}),
		}),
		[]byte("hello, world!"),
	)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(msg.ToBytes()))

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusInternalServerError {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingRedirectRequestGetAliasByPubKey(handler http.HandlerFunc) error {
	msg := layer3.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{}),
		}),
		[]byte("hello, world!"),
	)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(msg.ToBytes()))
	req.Header.Add(hls_settings.CHeaderPublicKey, asymmetric.NewPrivKey().GetPubKey().ToString())

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotExtended {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingRedirectRequestRedirect(handler http.HandlerFunc, pubKey asymmetric.IPubKey) error {
	msg := layer3.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{}),
		}),
		[]byte("hello, world!"),
	)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(msg.ToBytes()))
	req.Header.Add(hls_settings.CHeaderPublicKey, pubKey.ToString())

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadGateway {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingRedirectRequestOK(handler http.HandlerFunc, pubKey asymmetric.IPubKey) error {
	msg := layer3.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{}),
		}),
		[]byte("hello, world!"),
	)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(msg.ToBytes()))
	req.Header.Add(hls_settings.CHeaderPublicKey, pubKey.ToString())

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}
