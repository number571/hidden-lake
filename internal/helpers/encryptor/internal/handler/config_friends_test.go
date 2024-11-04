package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/app/config"
	hle_client "github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/client"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestHandleConfigFriendsAPI(t *testing.T) {
	t.Parallel()

	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 0)
	defer os.Remove(pathCfg)

	wcfg, service := testRunService(pathCfg, testutils.TgAddrs[39])
	defer service.Close()

	time.Sleep(100 * time.Millisecond)
	hleClient := hle_client.NewClient(
		hle_client.NewBuilder(),
		hle_client.NewRequester(
			"http://"+testutils.TgAddrs[39],
			&http.Client{Timeout: time.Second / 2},
			testNetworkMessageSettings(),
		),
	)

	aliasName := "test_name4"
	testGetFriends(t, hleClient, wcfg.GetConfig())
	testAddFriend(t, hleClient, aliasName)
	testDelFriend(t, hleClient, aliasName)
}

func testGetFriends(t *testing.T, client hle_client.IClient, cfg config.IConfig) {
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

func testAddFriend(t *testing.T, client hle_client.IClient, aliasName string) {
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

func testDelFriend(t *testing.T, client hle_client.IClient, aliasName string) {
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
