package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/service/pkg/app/config"
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
	"github.com/number571/hidden-lake/internal/service/pkg/settings"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestHandleFriendsAPI2(t *testing.T) {
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

	handler := HandleConfigFriendsAPI(newTsWrapper(true), httpLogger, newTsNode(true, true, true, true))
	if err := friendsAPIRequestOK(handler); err != nil {
		t.Error(err)
		return
	}
	if err := friendsAPIRequestPostOK(handler); err != nil {
		t.Error(err)
		return
	}
	if err := friendsAPIRequestDeleteOK(handler); err != nil {
		t.Error(err)
		return
	}

	if err := friendsAPIRequestNotFound(handler); err == nil {
		t.Error("request success with not found alias_name")
		return
	}
	if err := friendsAPIRequestPubKey(handler); err == nil {
		t.Error("request success with invalid pubkey")
		return
	}
	if err := friendsAPIRequestExist(handler); err == nil {
		t.Error("request success with exist alias_name")
		return
	}
	if err := friendsAPIRequestAliasName(handler); err == nil {
		t.Error("request success with invalid alias_name")
		return
	}
	if err := friendsAPIRequestDecode(handler); err == nil {
		t.Error("request success with invalid decode")
		return
	}
	if err := friendsAPIRequestMethod(handler); err == nil {
		t.Error("request success with invalid method")
		return
	}

	handlerx := HandleConfigFriendsAPI(newTsWrapper(false), httpLogger, newTsNode(true, true, true, true))
	if err := friendsAPIRequestPostOK(handlerx); err == nil {
		t.Error("request success with invalid update editor (post)")
		return
	}
	if err := friendsAPIRequestDeleteOK(handlerx); err == nil {
		t.Error("request success with invalid update editor (post)")
		return
	}
}

func friendsAPIRequestDeleteOK(handler http.HandlerFunc) error {
	newFriend := settings.SFriend{
		FAliasName: "abc",
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", bytes.NewBuffer(encoding.SerializeJSON(newFriend)))

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

func friendsAPIRequestPostOK(handler http.HandlerFunc) error {
	newFriend := settings.SFriend{
		FAliasName: "new_friend",
		FPublicKey: tgPrivKey3.GetPubKey().ToString(),
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(encoding.SerializeJSON(newFriend)))

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

func friendsAPIRequestNotFound(handler http.HandlerFunc) error {
	newFriend := settings.SFriend{
		FAliasName: "notfound",
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", bytes.NewBuffer(encoding.SerializeJSON(newFriend)))

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

func friendsAPIRequestOK(handler http.HandlerFunc) error {
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

func friendsAPIRequestPubKey(handler http.HandlerFunc) error {
	newFriend := settings.SFriend{
		FAliasName: "new_friend",
		FPublicKey: "abc",
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(encoding.SerializeJSON(newFriend)))

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

func friendsAPIRequestExist(handler http.HandlerFunc) error {
	newFriend := settings.SFriend{
		FAliasName: "abc",
		FPublicKey: tgPrivKey3.GetPubKey().ToString(),
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(encoding.SerializeJSON(newFriend)))

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

func friendsAPIRequestAliasName(handler http.HandlerFunc) error {
	newFriend := settings.SFriend{
		FAliasName: "",
		FPublicKey: tgPrivKey3.GetPubKey().ToString(),
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(encoding.SerializeJSON(newFriend)))

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

func friendsAPIRequestDecode(handler http.HandlerFunc) error {
	newFriend := settings.SFriend{
		FAliasName: "new_friend",
		FPublicKey: tgPrivKey3.GetPubKey().ToString(),
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bytes.Join(
		[][]byte{
			[]byte{1},
			encoding.SerializeJSON(newFriend),
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

func friendsAPIRequestMethod(handler http.HandlerFunc) error {
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

func TestHandleFriendsAPI(t *testing.T) {
	t.Parallel()

	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 1)
	pathDB := fmt.Sprintf(tcPathDBTemplate, 1)

	wcfg, node, _, cancel, srv := testAllCreate(pathCfg, pathDB, testutils.TgAddrs[7])
	defer testAllFree(node, cancel, srv, pathCfg, pathDB)

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			"http://"+testutils.TgAddrs[7],
			&http.Client{Timeout: time.Minute},
		),
	)

	aliasName := "test_name4"
	testGetFriends(t, client, wcfg.GetConfig())
	testAddFriend(t, client, aliasName)
	testDelFriend(t, client, aliasName)
}

func testGetFriends(t *testing.T, client hls_client.IClient, cfg config.IConfig) {
	friends, err := client.GetFriends(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if len(friends) != 2 {
		t.Error("length of friends != 2")
		return
	}

	for k, v := range friends {
		v1, ok := cfg.GetFriends()[k]
		if !ok {
			t.Errorf("undefined friend '%s'", k)
			return
		}
		if v.ToString() != v1.ToString() {
			t.Errorf("public keys not equals for '%s'", k)
			return
		}
	}
}

func testAddFriend(t *testing.T, client hls_client.IClient, aliasName string) {
	err := client.AddFriend(
		context.Background(),
		aliasName,
		tgPrivKey3.GetPubKey(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	friends, err := client.GetFriends(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if _, ok := friends[aliasName]; !ok {
		t.Errorf("undefined new public key by '%s'", aliasName)
		return
	}
}

func testDelFriend(t *testing.T, client hls_client.IClient, aliasName string) {
	err := client.DelFriend(context.Background(), aliasName)
	if err != nil {
		t.Error(err)
		return
	}

	friends, err := client.GetFriends(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if _, ok := friends[aliasName]; ok {
		t.Errorf("deleted public key exists for '%s'", aliasName)
		return
	}
}
