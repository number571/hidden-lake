package handler

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestHandleNetworkKeyAPI(t *testing.T) {
	t.Parallel()

	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 4)
	pathDB := fmt.Sprintf(tcPathDBTemplate, 4)

	_, node, _, cancel, srv := testAllCreate(pathCfg, pathDB, testutils.TgAddrs[14])
	defer testAllFree(node, cancel, srv, pathCfg, pathDB)

	client := hlk_client.NewClient(
		hlk_client.NewBuilder(),
		hlk_client.NewRequester(
			testutils.TgAddrs[14],
			&http.Client{Timeout: time.Minute},
		),
	)

	testGetNetworkKey(t, client, "test")
}

func testGetNetworkKey(t *testing.T, client hlk_client.IClient, networkKey string) {
	settings, err := client.GetSettings(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if settings.GetNetworkKey() != networkKey {
		t.Fatal("got network key != networkKey")
	}
}
