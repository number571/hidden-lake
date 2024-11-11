package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/logger"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/cache"
	hiddenlake "github.com/number571/hidden-lake"
	hlt_database "github.com/number571/hidden-lake/internal/helpers/traffic/internal/database"
	"github.com/number571/hidden-lake/internal/helpers/traffic/internal/storage"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestHandlePointerAPI2(t *testing.T) {
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

	storage := storage.NewMessageStorage(
		net_message.NewSettings(&net_message.SSettings{}),
		hlt_database.NewVoidKVDatabase(),
		cache.NewLRUCache(16),
	)

	handler := HandlePointerAPI(storage, httpLogger)
	if err := pointerAPIRequestOK(handler); err != nil {
		t.Error(err)
		return
	}

	if err := pointerAPIRequestMethod(handler); err == nil {
		t.Error("request success with invalid method")
		return
	}
}

func pointerAPIRequestOK(handler http.HandlerFunc) error {
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

func pointerAPIRequestMethod(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)

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

func TestHandlePointerAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[18]
	os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))

	srv, cancel, db, hltClient := testAllRun(addr)
	defer testAllFree(addr, srv, cancel, db)

	client := testNewClient()
	msg, err := client.EncryptMessage(
		client.GetPrivKey().GetPubKey(),
		payload.NewPayload64(0, []byte(tcBody)).ToBytes(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	netMsg := net_message.NewMessage(
		testNetworkMessageSettings(),
		payload.NewPayload32(hiddenlake.GSettings.FProtoMask.FService, msg),
	)
	if err := hltClient.PutMessage(context.Background(), netMsg); err != nil {
		t.Error(err)
		return
	}

	pointer, err := hltClient.GetPointer(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if pointer != 1 {
		t.Error("incorrect pointer")
		return
	}
}
