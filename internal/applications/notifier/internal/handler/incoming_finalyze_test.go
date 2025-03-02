package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/hidden-lake/internal/applications/notifier/pkg/app/config"
	"github.com/number571/hidden-lake/internal/utils/layer3"
	"github.com/number571/hidden-lake/internal/utils/msgdata"
)

func TestHandleIncomingFinalyzeHTTP(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	log := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	hlsClient := newTsHLSClient(true, true)
	handler := HandleIncomingFinalyzeHTTP(
		ctx,
		&config.SConfig{FSettings: &config.SConfigSettings{}},
		log,
		newTsDatabase(true, true),
		msgdata.NewMessageBroker(),
		hlsClient,
	)

	if err := incomingFinalyzeRequestMethod(handler); err != nil {
		t.Error(err)
		return
	}
	if err := incomingFinalyzeRequestLoadMessage(handler); err != nil {
		t.Error(err)
		return
	}
	if err := incomingFinalyzeRequestExtractMessage(handler); err != nil {
		t.Error(err)
		return
	}

	hlsClient2 := newTsHLSClient(false, true)
	handler2 := HandleIncomingFinalyzeHTTP(
		ctx,
		&config.SConfig{FSettings: &config.SConfigSettings{}},
		log,
		newTsDatabase(true, true),
		msgdata.NewMessageBroker(),
		hlsClient2,
	)
	if err := incomingFinalyzeRequestGetPubKey(handler2); err != nil {
		t.Error(err)
		return
	}

	database3 := newTsDatabase(false, true)
	database3.fSetHashWithoutOK = true
	handler3 := HandleIncomingFinalyzeHTTP(
		ctx,
		&config.SConfig{FSettings: &config.SConfigSettings{}},
		log,
		database3,
		msgdata.NewMessageBroker(),
		hlsClient,
	)
	if err := incomingFinalyzeRequestSetHash(handler3); err != nil {
		t.Error(err)
		return
	}

	// ...
}

func incomingFinalyzeRequestMethod(handler http.HandlerFunc) error {
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

func incomingFinalyzeRequestLoadMessage(handler http.HandlerFunc) error {
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

func incomingFinalyzeRequestExtractMessage(handler http.HandlerFunc) error {
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

func incomingFinalyzeRequestGetPubKey(handler http.HandlerFunc) error {
	msg := layer3.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{}),
		}),
		[]byte("hello"),
	)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(msg.ToBytes()))

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadGateway {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingFinalyzeRequestSetHash(handler http.HandlerFunc) error {
	msg := layer3.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{}),
		}),
		[]byte("hello"),
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
