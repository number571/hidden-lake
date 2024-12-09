package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/payload"
	hle_client "github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/client"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestHandleEncryptDecryptAPI(t *testing.T) {
	t.Parallel()

	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 2)
	defer os.Remove(pathCfg)

	_, service := testRunService(pathCfg, testutils.TgAddrs[33])
	defer service.Close()

	time.Sleep(100 * time.Millisecond)
	hleClient := hle_client.NewClient(
		hle_client.NewBuilder(),
		hle_client.NewRequester(
			testutils.TgAddrs[33],
			&http.Client{Timeout: time.Second / 2},
			testNetworkMessageSettings(),
		),
	)

	// same private key in the HLE
	data := []byte("hello, world!")
	netMsg, err := hleClient.EncryptMessage(
		context.Background(),
		"test_recvr",
		payload.NewPayload64(1, data),
	)
	if err != nil {
		t.Error(err)
		return
	}

	aliasName, getPld, err := hleClient.DecryptMessage(context.Background(), netMsg)
	if err != nil {
		t.Error(err)
		return
	}

	if aliasName != "test_recvr" {
		t.Error("got invalid alias_name")
		return
	}

	if !bytes.Equal(getPld.GetBody(), data) {
		t.Error("got invalid data")
		return
	}
}
