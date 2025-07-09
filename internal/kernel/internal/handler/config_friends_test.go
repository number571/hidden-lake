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
	"github.com/number571/hidden-lake/internal/kernel/pkg/app/config"
	hls_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	"github.com/number571/hidden-lake/internal/kernel/pkg/settings"
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

	handler := HandleConfigFriendsAPI(newTsWrapper(true), httpLogger, newTsNode(true, true, true))
	if err := friendsAPIRequestOK(handler); err != nil {
		t.Fatal(err)
	}
	if err := friendsAPIRequestPostOK(handler); err != nil {
		t.Fatal(err)
	}
	if err := friendsAPIRequestDeleteOK(handler); err != nil {
		t.Fatal(err)
	}

	if err := friendsAPIRequestNotFound(handler); err == nil {
		t.Fatal("request success with not found alias_name")
	}
	if err := friendsAPIRequestPubKey(handler); err == nil {
		t.Fatal("request success with invalid pubkey")
	}
	if err := friendsAPIRequestExist(handler); err == nil {
		t.Fatal("request success with exist alias_name")
	}
	if err := friendsAPIRequestAliasName(handler); err == nil {
		t.Fatal("request success with invalid alias_name")
	}
	if err := friendsAPIRequestDecode(handler); err == nil {
		t.Fatal("request success with invalid decode")
	}
	if err := friendsAPIRequestMethod(handler); err == nil {
		t.Fatal("request success with invalid method")
	}

	handlerx := HandleConfigFriendsAPI(newTsWrapper(false), httpLogger, newTsNode(true, true, true))
	if err := friendsAPIRequestPostOK(handlerx); err == nil {
		t.Fatal("request success with invalid update editor (post)")
	}
	if err := friendsAPIRequestDeleteOK(handlerx); err == nil {
		t.Fatal("request success with invalid update editor (post)")
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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
	defer func() { _ = res.Body.Close() }()

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
			testutils.TgAddrs[7],
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
		t.Fatal(err)
	}

	if len(friends) != 2 {
		t.Fatal("length of friends != 2")
	}

	for k, v := range friends {
		v1, ok := cfg.GetFriends()[k]
		if !ok {
			t.Fatalf("undefined friend '%s'", k)
		}
		if v.ToString() != v1.ToString() {
			t.Fatalf("public keys not equals for '%s'", k)
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
		t.Fatal(err)
	}

	friends, err := client.GetFriends(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := friends[aliasName]; !ok {
		t.Fatalf("undefined new public key by '%s'", aliasName)
	}
}

func testDelFriend(t *testing.T, client hls_client.IClient, aliasName string) {
	err := client.DelFriend(context.Background(), aliasName)
	if err != nil {
		t.Fatal(err)
	}

	friends, err := client.GetFriends(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := friends[aliasName]; ok {
		t.Fatalf("deleted public key exists for '%s'", aliasName)
	}
}
