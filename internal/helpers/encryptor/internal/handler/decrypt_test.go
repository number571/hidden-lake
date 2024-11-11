package handler

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	hiddenlake "github.com/number571/hidden-lake"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func TestHandleMessageDecryptAPI(t *testing.T) {
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

	clientSender := client.NewClient(tgPrivKey2, 8192)
	encMsg, err := clientSender.EncryptMessage(
		tgPrivKey1.GetPubKey(),
		payload.NewPayload64(1, []byte("hello, world!")).ToBytes(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	netMsg := net_message.NewMessage(
		net_message.NewConstructSettings(&net_message.SConstructSettings{
			FSettings: net_message.NewSettings(&net_message.SSettings{}),
		}),
		payload.NewPayload32(hiddenlake.GSettings.FProtoMask.FService, encMsg),
	).ToString()

	mapKeys := asymmetric.NewMapPubKeys()
	mapKeys.SetPubKey(tgPrivKey2.GetPubKey())

	clientReceiver := client.NewClient(tgPrivKey1, 8192)
	handler := HandleMessageDecryptAPI(&tsConfig{}, httpLogger, clientReceiver, mapKeys)
	if err := decryptAPIRequestOK(handler, netMsg); err != nil {
		t.Error(err)
		return
	}

	if err := decryptAPIRequestDecode(handler, netMsg); err == nil {
		t.Error("request success with invalid method")
		return
	}
	if err := decryptAPIRequestMethod(handler); err == nil {
		t.Error("request success with invalid method")
		return
	}

	netMsgHeadCustom := net_message.NewMessage(
		net_message.NewConstructSettings(&net_message.SConstructSettings{
			FSettings: net_message.NewSettings(&net_message.SSettings{}),
		}),
		payload.NewPayload32(1, encMsg),
	).ToString()
	if err := decryptAPIRequestOK(handler, netMsgHeadCustom); err == nil {
		t.Error("request success with invalid head (custom)")
		return
	}

	encMsgNotFayload, err := clientSender.EncryptMessage(
		tgPrivKey1.GetPubKey(),
		[]byte{123},
	)
	if err != nil {
		t.Error(err)
		return
	}
	netMsgNotPayload := net_message.NewMessage(
		net_message.NewConstructSettings(&net_message.SConstructSettings{
			FSettings: net_message.NewSettings(&net_message.SSettings{}),
		}),
		payload.NewPayload32(hiddenlake.GSettings.FProtoMask.FService, encMsgNotFayload),
	).ToString()
	if err := decryptAPIRequestOK(handler, netMsgNotPayload); err == nil {
		t.Error("request success with invalid decrypted payload")
		return
	}

	handlerx := HandleMessageDecryptAPI(&tsConfig{}, httpLogger, clientReceiver, asymmetric.NewMapPubKeys())
	if err := decryptAPIRequestOK(handlerx, netMsg); err == nil {
		t.Error("request success with invalid decrypt")
		return
	}

}

func decryptAPIRequestOK(handler http.HandlerFunc, netMsg string) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(netMsg)))

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

func decryptAPIRequestDecode(handler http.HandlerFunc, netMsg string) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bytes.Join(
		[][]byte{
			[]byte{1},
			[]byte(netMsg),
		},
		[]byte{},
	)))

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

func decryptAPIRequestMethod(handler http.HandlerFunc) error {
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
