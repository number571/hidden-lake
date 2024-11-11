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

func TestHandleMessageAPI2(t *testing.T) {
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

	client := client.NewClient(asymmetric.NewPrivKey(), tcMessageSize)
	msg, err := client.EncryptMessage(
		client.GetPrivKey().GetPubKey(),
		payload.NewPayload64(0, []byte("hello")).ToBytes(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	netMsg := net_message.NewMessage(
		net_message.NewConstructSettings(&net_message.SConstructSettings{
			FSettings: net_message.NewSettings(&net_message.SSettings{}),
		}),
		payload.NewPayload32(hiddenlake.GSettings.FProtoMask.FService, msg),
	)
	if err := storage.Push(netMsg); err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()
	handler := HandleMessageAPI(ctx, &tsConfig{}, storage, httpLogger, httpLogger, newNetworkNode())
	if err := messageAPIRequestOK(handler, encoding.HexEncode(netMsg.GetHash())); err != nil {
		t.Error(err)
		return
	}
	if err := messageAPIRequestPostOK(handler, netMsg.ToString()); err != nil {
		t.Error(err)
		return
	}
	if err := messageAPIRequestPostOK(handler, netMsg.ToString()); err == nil {
		t.Error("request success with duplicate message")
		return
	}

	netMsgCustom := net_message.NewMessage(
		net_message.NewConstructSettings(&net_message.SConstructSettings{
			FSettings: net_message.NewSettings(&net_message.SSettings{}),
		}),
		payload.NewPayload32(1, []byte("hello")),
	)

	if err := messageAPIRequestPostOK(handler, netMsgCustom.ToString()); err == nil {
		t.Error("request success with invalid message (custom)")
		return
	}
	if err := messageAPIRequestPostOK(handler, "abc"); err == nil {
		t.Error("request success with invalid message")
		return
	}
	if err := messageAPIRequestOK(handler, "hello!"); err == nil {
		t.Error("request success with invalid hash")
		return
	}
	if err := messageAPIRequestOK(handler, "eax"); err == nil {
		t.Error("request success with load from hash")
		return
	}
	if err := messageAPIRequestMethod(handler); err == nil {
		t.Error("request success with invalid method")
		return
	}
}

func messageAPIRequestOK(handler http.HandlerFunc, hash string) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?hash="+hash, nil)

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

func messageAPIRequestPostOK(handler http.HandlerFunc, msg string) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(msg)))

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

func messageAPIRequestMethod(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", nil)

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

func TestHandleMessageAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[20]
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

	strHash := encoding.HexEncode(netMsg.GetHash())
	gotNetMsg, err := hltClient.GetMessage(context.Background(), strHash)
	if err != nil {
		t.Error(err)
		return
	}

	gotPubKey, decMsg, err := client.DecryptMessage(
		asymmetric.NewMapPubKeys(client.GetPrivKey().GetPubKey()),
		gotNetMsg.GetPayload().GetBody(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(gotPubKey.ToBytes(), client.GetPrivKey().GetPubKey().ToBytes()) {
		t.Error("invalid public keys")
		return
	}

	gotPld := payload.LoadPayload64(decMsg)
	if string(gotPld.GetBody()) != tcBody {
		t.Error(err)
		return
	}
}
