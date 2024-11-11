package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
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

func TestHandleHashesAPI2(t *testing.T) {
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

	handler := HandleHashesAPI(storage, httpLogger)
	if err := hashesAPIRequestOK(handler, "0"); err == nil {
		t.Error("request success with none hashes")
		return
	}

	err := storage.Push(
		net_message.NewMessage(
			net_message.NewConstructSettings(&net_message.SConstructSettings{
				FSettings: net_message.NewSettings(&net_message.SSettings{}),
			}),
			payload.NewPayload32(1, []byte("hello")),
		),
	)
	if err != nil {
		t.Error(err)
		return
	}

	if err := hashesAPIRequestOK(handler, "0"); err != nil {
		t.Error(err)
		return
	}

	if err := hashesAPIRequestOK(handler, "abc"); err == nil {
		t.Error("request success with invalid id")
		return
	}
	if err := hashesAPIRequestMethod(handler); err == nil {
		t.Error("request success with invalid method")
		return
	}
}

func hashesAPIRequestOK(handler http.HandlerFunc, id string) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?id="+id, nil)

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

func hashesAPIRequestMethod(handler http.HandlerFunc) error {
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

func TestHandleHashesAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[19]
	os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))

	srv, cancel, wDB, hltClient := testAllRun(addr)
	defer testAllFree(addr, srv, cancel, wDB)

	privKey := asymmetric.NewPrivKey()
	pubKey := privKey.GetPubKey()

	client := client.NewClient(privKey, tcMessageSize)
	msg, err := client.EncryptMessage(
		pubKey,
		payload.NewPayload64(0, []byte("hello")).ToBytes(),
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

	hash, err := hltClient.GetHash(context.Background(), 0)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(encoding.HexDecode(hash), netMsg.GetHash()) {
		t.Error("hashes not equals")
		return
	}
}
