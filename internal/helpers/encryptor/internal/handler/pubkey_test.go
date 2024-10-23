package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	hle_client "github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/client"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestHandlePubKeyAPI(t *testing.T) {
	t.Parallel()

	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 3)
	defer os.Remove(pathCfg)

	_, service := testRunService(pathCfg, testutils.TgAddrs[31])
	defer service.Close()

	time.Sleep(100 * time.Millisecond)
	hleClient := hle_client.NewClient(
		hle_client.NewBuilder(),
		hle_client.NewRequester(
			"http://"+testutils.TgAddrs[31],
			&http.Client{Timeout: time.Second / 2},
			testNetworkMessageSettings(),
		),
	)

	gotPubKey, err := hleClient.GetPubKey(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	pubKey := tgPrivKey1.GetPubKey()
	if pubKey.ToString() != gotPubKey.ToString() {
		t.Error("public keys not equals")
		return
	}
}
