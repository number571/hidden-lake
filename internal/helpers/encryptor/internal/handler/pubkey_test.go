package handler

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	testutils "github.com/number571/go-peer/test/utils"
	hle_client "github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/client"
)

func TestHandlePubKeyAPI(t *testing.T) {
	t.Parallel()

	service := testRunService(testutils.TgAddrs[49])
	defer service.Close()

	time.Sleep(100 * time.Millisecond)
	hleClient := hle_client.NewClient(
		hle_client.NewRequester(
			"http://"+testutils.TgAddrs[49],
			&http.Client{Timeout: time.Second / 2},
			testNetworkMessageSettings(),
		),
	)

	gotPubKey, err := hleClient.GetPubKey(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	pubKey := asymmetric.LoadRSAPrivKey(testutils.TcPrivKey1024).GetPubKey()
	if !bytes.Equal(gotPubKey.ToBytes(), pubKey.ToBytes()) {
		t.Error("public keys not equals")
		return
	}
}
