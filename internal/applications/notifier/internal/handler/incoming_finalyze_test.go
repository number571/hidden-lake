package handler

// import (
// 	"bytes"
// 	"context"
// 	"errors"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/number571/go-peer/pkg/crypto/asymmetric"
// 	"github.com/number571/go-peer/pkg/logger"
// 	"github.com/number571/hidden-lake/internal/applications/notifier/internal/database"
// 	"github.com/number571/hidden-lake/internal/applications/notifier/internal/msgbroker"
// 	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
// 	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
// )

// func TestHandleIncomingPushHTTP(t *testing.T) {
// 	t.Parallel()

// 	logging, err := std_logger.LoadLogging([]string{})
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	httpLogger := std_logger.NewStdLogger(
// 		logging,
// 		func(_ logger.ILogArg) string {
// 			return ""
// 		},
// 	)

// 	ctx := context.Background()
// 	msgBroker := msgdata.NewMessageBroker()
// 	handler := HandleIncomingPushHTTP(ctx, httpLogger, newTsDatabase(true, true), msgBroker, newTsHLSClient(true, true))

// 	if err := incomingPushRequestOK(handler); err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	if err := incomingPushRequestMethod(handler); err == nil {
// 		t.Error("request success with invalid method")
// 		return
// 	}
// 	if err := incomingPushRequestPubKey(handler); err == nil {
// 		t.Error("request success with invalid pubkey")
// 		return
// 	}
// 	if err := incomingPushRequestMessage(handler); err == nil {
// 		t.Error("request success with invalid message")
// 		return
// 	}

// 	handlerx := HandleIncomingPushHTTP(ctx, httpLogger, newTsDatabase(true, true), msgBroker, newTsHLSClient(false, true))
// 	if err := incomingPushRequestOK(handlerx); err == nil {
// 		t.Error("request success with invalid my pubkey")
// 		return
// 	}
// 	handlery := HandleIncomingPushHTTP(ctx, httpLogger, newTsDatabase(false, true), msgBroker, newTsHLSClient(true, true))
// 	if err := incomingPushRequestOK(handlery); err == nil {
// 		t.Error("request success with invalid push message")
// 		return
// 	}
// }

// func incomingPushRequestOK(handler http.HandlerFunc) error {
// 	w := httptest.NewRecorder()
// 	req := httptest.NewRequest(http.MethodPost, "/push", bytes.NewBuffer(wrapText("hello, world!")))
// 	req.Header.Set(hls_settings.CHeaderPublicKey, asymmetric.NewPrivKey().GetPubKey().ToString())

// 	handler(w, req)
// 	res := w.Result()
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		return errors.New("bad status code") // nolint: err113
// 	}

// 	if _, err := io.ReadAll(res.Body); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func incomingPushRequestMessage(handler http.HandlerFunc) error {
// 	w := httptest.NewRecorder()
// 	req := httptest.NewRequest(http.MethodPost, "/push", nil)
// 	req.Header.Set(hls_settings.CHeaderPublicKey, asymmetric.NewPrivKey().GetPubKey().ToString())

// 	handler(w, req)
// 	res := w.Result()
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		return errors.New("bad status code") // nolint: err113
// 	}

// 	if _, err := io.ReadAll(res.Body); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func incomingPushRequestPubKey(handler http.HandlerFunc) error {
// 	w := httptest.NewRecorder()
// 	req := httptest.NewRequest(http.MethodPost, "/push", bytes.NewBuffer(wrapText("hello, world!")))

// 	handler(w, req)
// 	res := w.Result()
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		return errors.New("bad status code") // nolint: err113
// 	}

// 	if _, err := io.ReadAll(res.Body); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func incomingPushRequestMethod(handler http.HandlerFunc) error {
// 	w := httptest.NewRecorder()
// 	req := httptest.NewRequest(http.MethodGet, "/push", nil)

// 	handler(w, req)
// 	res := w.Result()
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		return errors.New("bad status code") // nolint: err113
// 	}

// 	return nil
// }

// type tsDatabase struct {
// 	fPushOK bool
// 	fLoadOK bool
// 	fMsg    database.IMessage
// }

// func newTsDatabase(pPushOK, pLoadOK bool) *tsDatabase {
// 	return &tsDatabase{
// 		fPushOK: pPushOK,
// 		fLoadOK: pLoadOK,
// 	}
// }

// func (p *tsDatabase) Close() error { return nil }

// func (p *tsDatabase) Size(database.IRelation) uint64 {
// 	if p.fMsg == nil {
// 		return 0
// 	}
// 	return 1
// }

// func (p *tsDatabase) Push(_ database.IRelation, pM database.IMessage) error {
// 	if !p.fPushOK {
// 		return errors.New("some error") // nolint: err113
// 	}
// 	p.fMsg = pM
// 	return nil
// }

// func (p *tsDatabase) Load(database.IRelation, uint64, uint64) ([]database.IMessage, error) {
// 	if !p.fLoadOK {
// 		return nil, errors.New("some error") // nolint: err113
// 	}
// 	if p.fMsg == nil {
// 		return nil, nil
// 	}
// 	return []database.IMessage{p.fMsg}, nil
// }
