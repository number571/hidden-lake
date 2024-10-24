package handler

import (
	"context"
	"fmt"
	"os"
	"testing"

	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestHandlePointerAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[18]
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
		payload.NewPayload32(hls_settings.CNetworkMask, msg),
	)
	if err := hltClient.PutMessage(context.Background(), netMsg); err != nil {
		t.Error(err)
		return
	}

	pointer, err := hltClient.GetPointer(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if pointer != 1 {
		t.Error("incorrect pointer")
		return
	}
}
