package handler

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/encoding"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestHandleMessageAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[20]
	os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))

	srv, cancel, db, hltClient := testAllRun(addr)
	defer testAllFree(addr, srv, cancel, db)

	client := testNewClient()
	msg, err := client.EncryptMessage(
		client.GetPrivKey().GetKEMPrivKey().GetPubKey(),
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

	strHash := encoding.HexEncode(netMsg.GetHash())
	gotNetMsg, err := hltClient.GetMessage(context.Background(), strHash)
	if err != nil {
		t.Error(err)
		return
	}

	gotPubKey, decMsg, err := client.DecryptMessage(gotNetMsg.GetPayload().GetBody())
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
