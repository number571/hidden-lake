package handler

import (
	"bytes"
	"context"
	"errors"
	"hash/crc32"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
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

	database3X := newTsDatabase(false, true)
	database3X.fSetHashWithoutOK = true
	database3X.fFailedAfterCounter = 1
	handler3X := HandleIncomingFinalyzeHTTP(
		ctx,
		&config.SConfig{
			FSettings: &config.SConfigSettings{},
			FChannels: []string{"111", "key", "222"},
		},
		log,
		database3X,
		msgdata.NewMessageBroker(),
		hlsClient,
	)
	if err := incomingFinalyzeRequestDecryptedSetHash(handler3X); err != nil {
		t.Error(err)
		return
	}

	hlsClient4X := newTsHLSClient(true, true)
	hlsClient4X.fFriends = map[string]asymmetric.IPubKey{
		"111": asymmetric.NewPrivKey().GetPubKey(),
		"222": asymmetric.NewPrivKey().GetPubKey(),
		"333": asymmetric.NewPrivKey().GetPubKey(),
	}
	handler4X := HandleIncomingFinalyzeHTTP(
		ctx,
		&config.SConfig{
			FSettings: &config.SConfigSettings{},
			FChannels: []string{"111", "key", "222"},
		},
		log,
		newTsDatabase(true, true),
		msgdata.NewMessageBroker(),
		hlsClient4X,
	)
	if err := incomingFinalyzeRequestDecryptedExtractMessage(handler4X); err != nil {
		t.Error(err)
		return
	}
	if err := incomingFinalyzeRequestDecryptedGetMessage(handler4X); err != nil {
		t.Error(err)
		return
	}

	database5X := newTsDatabase(false, true)
	database5X.fFailedAfterCounter = 2
	handler5X := HandleIncomingFinalyzeHTTP(
		ctx,
		&config.SConfig{
			FSettings: &config.SConfigSettings{},
			FChannels: []string{"111", "key", "222"},
		},
		log,
		database5X,
		msgdata.NewMessageBroker(),
		hlsClient,
	)
	if err := incomingFinalyzeRequestDecryptedPush(handler5X); err != nil {
		t.Error(err)
		return
	}
	if err := incomingFinalyzeRequestDecryptedOK(handler4X); err != nil {
		t.Error(err)
		return
	}

	if err := incomingFinalyzeRequestGetFriends(handler); err != nil {
		t.Error(err)
		return
	}

	hlsClient4 := newTsHLSClient(true, false)
	hlsClient4.fFriends = map[string]asymmetric.IPubKey{
		"111": asymmetric.NewPrivKey().GetPubKey(),
		"222": asymmetric.NewPrivKey().GetPubKey(),
		"333": asymmetric.NewPrivKey().GetPubKey(),
	}
	handler4 := HandleIncomingFinalyzeHTTP(
		ctx,
		&config.SConfig{FSettings: &config.SConfigSettings{}},
		log,
		newTsDatabase(true, true),
		msgdata.NewMessageBroker(),
		hlsClient4,
	)

	if err := incomingFinalyzeRequestFinalyze(handler4); err != nil {
		t.Error(err)
		return
	}

	hlsClient5 := newTsHLSClient(true, true)
	hlsClient5.fFriends = map[string]asymmetric.IPubKey{
		"111": asymmetric.NewPrivKey().GetPubKey(),
		"222": asymmetric.NewPrivKey().GetPubKey(),
		"333": asymmetric.NewPrivKey().GetPubKey(),
	}
	handler5 := HandleIncomingFinalyzeHTTP(
		ctx,
		&config.SConfig{FSettings: &config.SConfigSettings{}},
		log,
		newTsDatabase(true, true),
		msgdata.NewMessageBroker(),
		hlsClient5,
	)

	if err := incomingFinalyzeRequestOK(handler5); err != nil {
		t.Error(err)
		return
	}
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

	if res.StatusCode != http.StatusNotExtended {
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

func incomingFinalyzeRequestDecryptedSetHash(handler http.HandlerFunc) error {
	msg := layer3.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{
				FNetworkKey: "key",
			}),
		}),
		[]byte("hello"),
	)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(msg.ToBytes()))

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusServiceUnavailable {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingFinalyzeRequestDecryptedExtractMessage(handler http.HandlerFunc) error {
	msg := layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{}),
		}),
		newPayload(
			layer1.NewMessage(
				layer1.NewConstructSettings(&layer1.SConstructSettings{
					FSettings: layer1.NewSettings(&layer1.SSettings{
						FNetworkKey: "key",
					}),
				}),
				payload.NewPayload32(1, []byte("hello")),
			).ToBytes(),
		),
	)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(msg.ToBytes()))

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnsupportedMediaType {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingFinalyzeRequestDecryptedGetMessage(handler http.HandlerFunc) error {
	msg := layer3.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{
				FNetworkKey: "key",
			}),
		}),
		[]byte("hello"),
	)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(msg.ToBytes()))

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusTeapot {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingFinalyzeRequestDecryptedPush(handler http.HandlerFunc) error {
	const cIsText = 1
	msg := layer3.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{
				FNetworkKey: "key",
			}),
		}),
		bytes.Join([][]byte{
			{cIsText},
			[]byte("hello"),
		}, []byte{}),
	)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(msg.ToBytes()))

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusInsufficientStorage {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingFinalyzeRequestDecryptedOK(handler http.HandlerFunc) error {
	const cIsText = 1
	msg := layer3.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{
				FNetworkKey: "key",
			}),
		}),
		bytes.Join([][]byte{
			{cIsText},
			[]byte("hello"),
		}, []byte{}),
	)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(msg.ToBytes()))

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func newPayload(pBody []byte) payload.IPayload32 {
	const cSaltSize = 32
	saltBody := bytes.Join(
		[][]byte{random.NewRandom().GetBytes(cSaltSize), pBody},
		[]byte{},
	)
	return payload.NewPayload32(
		crc32.Checksum(saltBody, crc32.IEEETable),
		saltBody,
	)
}

func incomingFinalyzeRequestGetFriends(handler http.HandlerFunc) error {
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

	if res.StatusCode != http.StatusNotFound {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func incomingFinalyzeRequestFinalyze(handler http.HandlerFunc) error {
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

func incomingFinalyzeRequestOK(handler http.HandlerFunc) error {
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

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}
