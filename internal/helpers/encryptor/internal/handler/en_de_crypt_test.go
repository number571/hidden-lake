package handler

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/payload"
	hle_client "github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/client"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestHandleEncryptDecryptAPI(t *testing.T) {
	t.Parallel()

	service := testRunService(testutils.TgAddrs[33])
	defer service.Close()

	time.Sleep(100 * time.Millisecond)
	hleClient := hle_client.NewClient(
		hle_client.NewRequester(
			"http://"+testutils.TgAddrs[33],
			&http.Client{Timeout: time.Second / 2},
			testNetworkMessageSettings(),
		),
	)

	// same private key in the HLE
	pubKey := tgPrivKey.GetPubKey()
	data := []byte("hello, world!")

	netMsg, err := hleClient.EncryptMessage(
		context.Background(),
		pubKey.GetKEncPubKey(),
		payload.NewPayload64(1, data),
	)
	if err != nil {
		t.Error(err)
		return
	}

	_, getPld, err := hleClient.DecryptMessage(context.Background(), netMsg)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(getPld.GetBody(), data) {
		t.Error("got invalid data")
		return
	}
}
