package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/messenger/internal/message"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	"github.com/number571/hidden-lake/pkg/api/services/messenger/client/dto"
)

func TestHandleChatSubscribeAPI(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
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

	msg := dto.NewMessage(true, "hello, world!", time.Now())
	msgBroker := message.NewMessageBroker()
	go func() {
		time.Sleep(100 * time.Millisecond)
		msgBroker.Produce("abc", msg)
	}()

	handlerX := HandleChatSubscribeAPI(ctx, httpLogger, msgBroker)
	if err := chatSubscribeRequestOK(handlerX, msg); err != nil {
		t.Fatal(err)
	}
	if err := chatSubscribeRequestInvalidMethod(handlerX); err != nil {
		t.Fatal(err)
	}

	// used for check got message from another friend
	msgBroker.Produce("qwerty", msg)

	chCtx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	handlerY := HandleChatSubscribeAPI(chCtx, httpLogger, msgBroker)
	if err := chatSubscribeRequestNoContent(handlerY); err != nil {
		t.Fatal(err)
	}
}

func chatSubscribeRequestOK(handler http.HandlerFunc, msgSend dto.IMessage) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?friend=abc&sid=111", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	msg, err := dto.LoadMessage(body)
	if err != nil {
		return err
	}
	if msg.ToString() != msgSend.ToString() {
		return errors.New("invalid message") // nolint: err113
	}

	return nil
}

func chatSubscribeRequestNoContent(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?friend=abc&sid=111", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusNoContent {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

func chatSubscribeRequestInvalidMethod(handler http.HandlerFunc) error {
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
