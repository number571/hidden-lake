package handler

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/utils"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

func TestHandleHashesAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[19]
	os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))

	srv, cancel, wDB, hltClient := testAllRun(addr)
	defer testAllFree(addr, srv, cancel, wDB)

	privKey := asymmetric.NewPrivKeyChain(
		asymmetric.NewKEncPrivKey(),
		asymmetric.NewSignPrivKey(),
	)
	pubKey := privKey.GetPubKeyChain()

	client := client.NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: testutils.TCMessageSize,
			FEncKeySizeBytes:  asymmetric.CKEncSize,
		}),
		privKey,
	)

	msg, err := client.EncryptMessage(
		pubKey.GetKEncPubKey(),
		payload.NewPayload64(0, []byte("hello")).ToBytes(),
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

	hash, err := hltClient.GetHash(context.Background(), 0)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(encoding.HexDecode(hash), netMsg.GetHash()) {
		t.Error("hashes not equals")
		return
	}
}
