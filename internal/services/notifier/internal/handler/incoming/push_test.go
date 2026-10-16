package incoming

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/number571/go-peer/pkg/logger"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	"github.com/number571/hidden-lake/internal/services/notifier/internal/database"
	"github.com/number571/hidden-lake/internal/services/notifier/pkg/app/config"
	"github.com/number571/hidden-lake/internal/utils/broker"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	message_dto "github.com/number571/hidden-lake/pkg/api/services/notifier/client/dto"
)

func TestHandleIncomingMsgboxHTTP(t *testing.T) {
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

	cfg := &config.SConfig{
		FSettings: &config.SConfigSettings{},
	}

	ctx := context.Background()
	msgBroker := broker.NewDataBroker(256, 256)
	handler := HandleIncomingMsgboxHTTP(ctx, httpLogger, cfg, newTsDatabase(true, true), msgBroker)

	if err := incomingPushRequestOK(handler); err != nil {
		t.Fatal(err)
	}

	if err := incomingPushRequestMethod(handler); err == nil {
		t.Fatal("request success with invalid method")
	}
	if err := incomingPushRequestPubKey(handler); err == nil {
		t.Fatal("request success with invalid pubkey")
	}
	if err := incomingPushRequestMessage(handler); err == nil {
		t.Fatal("request success with invalid message")
	}

	handlery := HandleIncomingMsgboxHTTP(ctx, httpLogger, cfg, newTsDatabase(false, true), msgBroker)
	if err := incomingPushRequestOK(handlery); err == nil {
		t.Fatal("request success with invalid push message")
	}
}

func incomingPushRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()

	path := fmt.Sprintf("/msgbox?uid=%s", message_dto.NewMessageUID()) // nolint:perfsprint
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBuffer([]byte("hello, world!")))

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

func incomingPushRequestMessage(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/msgbox", nil)
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

func incomingPushRequestPubKey(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/msgbox", bytes.NewBuffer([]byte("hello, world!")))

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

func incomingPushRequestMethod(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/msgbox", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	return nil
}

type tsDatabase struct {
	fPushOK bool
	fLoadOK bool
	fMsg    message_dto.IMessage
}

func newTsDatabase(pPushOK, pLoadOK bool) *tsDatabase {
	return &tsDatabase{
		fPushOK: pPushOK,
		fLoadOK: pLoadOK,
	}
}

func (p *tsDatabase) Close() error { return nil }

func (p *tsDatabase) Size(string) uint64 {
	if p.fMsg == nil {
		return 0
	}
	return 1
}

func (p *tsDatabase) Push(_ string, pM message_dto.IMessage) (uint64, error) {
	if !p.fPushOK {
		return 0, errors.New("some error") // nolint: err113
	}
	p.fMsg = pM
	return 0, nil
}

func (p *tsDatabase) Recv(string, string) error {
	return nil
}

func (p *tsDatabase) LoadByUID(string, string) (message_dto.IMessage, error) {
	if !p.fLoadOK {
		return nil, errors.New("some error") // nolint: err113
	}
	if p.fMsg == nil {
		return nil, database.ErrNotFoundMessage
	}
	return p.fMsg, nil
}

func (p *tsDatabase) LoadByEnum(string, uint64) (message_dto.IMessage, error) {
	if !p.fLoadOK {
		return nil, errors.New("some error") // nolint: err113
	}
	if p.fMsg == nil {
		return nil, database.ErrNotFoundMessage
	}
	return p.fMsg, nil
}
