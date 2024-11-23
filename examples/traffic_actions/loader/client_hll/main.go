package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/hidden-lake/build"
	hll_client "github.com/number571/hidden-lake/internal/helpers/loader/pkg/client"
	hlt_client "github.com/number571/hidden-lake/internal/helpers/traffic/pkg/client"
)

const (
	hltHost1 = "localhost:7582"
	hltHost2 = "localhost:8582"
	hllHost  = "localhost:6561"
)

const (
	messageSize = (8 << 10) // 8KiB
	networkKey  = "some-network-key"
	workSize    = 10
	keySize     = 1024
)

const (
	messageCount = 64
)

var (
	privKey      asymmetric.IPrivKey
	pushedHashes = make([][]byte, 0, messageCount)
)

func init() {
	readPrivKey, err := os.ReadFile("priv_node1.key")
	if err != nil {
		panic(err)
	}
	privKey = asymmetric.LoadPrivKey(string(readPrivKey))
}

func main() {
	netMsgSettings := net_message.NewConstructSettings(&net_message.SConstructSettings{
		FSettings: net_message.NewSettings(&net_message.SSettings{
			FNetworkKey:   networkKey,
			FWorkSizeBits: workSize,
		}),
	})

	ctx := context.Background()
	if err := pushMessages(ctx, netMsgSettings, messageSize); err != nil {
		panic(err)
	}

	hllClient := hll_client.NewClient(
		hll_client.NewRequester(
			"http://"+hllHost,
			&http.Client{Timeout: time.Minute / 2},
		),
	)

	if err := hllClient.RunTransfer(ctx); err != nil {
		panic(err)
	}

	time.Sleep(time.Second)

	if err := checkMessages(ctx, netMsgSettings.GetSettings(), messageSize); err != nil {
		panic(err)
	}

	fmt.Println("messages have been successfully transported")
}

func pushMessages(ctx context.Context, netMsgSettings net_message.IConstructSettings, msgSize uint64) error {
	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			"http://"+hltHost1,
			&http.Client{Timeout: time.Minute / 2},
			netMsgSettings.GetSettings(),
		),
	)

	client := client.NewClient(privKey, msgSize)

	for i := 0; i < messageCount; i++ {
		msg, err := client.EncryptMessage(
			client.GetPrivKey().GetPubKey(), // self encrypt
			payload.NewPayload64(uint64(i), []byte("hello, world!")).ToBytes(),
		)
		if err != nil {
			return err
		}

		netMsg := net_message.NewMessage(
			netMsgSettings,
			payload.NewPayload32(build.GSettings.FProtoMask.FNetwork, msg),
		)
		if err := hltClient.PutMessage(ctx, netMsg); err != nil {
			return err
		}

		pushedHashes = append(pushedHashes, netMsg.GetHash())
	}

	return nil
}

func checkMessages(ctx context.Context, netMsgSettings net_message.ISettings, msgSize uint64) error {
	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			"http://"+hltHost2,
			&http.Client{Timeout: time.Minute / 2},
			netMsgSettings,
		),
	)

	client := client.NewClient(privKey, msgSize)

	hashes := make([]string, 0, messageCount)
	for i := uint64(0); ; i++ {
		hash, err := hltClient.GetHash(ctx, i)
		if err != nil {
			break
		}
		hashes = append(hashes, hash)
	}

	for _, ph := range pushedHashes {
		if !hashIsExist(ph, hashes) {
			return errors.New("hash not found")
		}
	}

	for _, h := range hashes {
		netMsg, err := hltClient.GetMessage(ctx, h)
		if err != nil {
			return err
		}

		if netMsg.GetPayload().GetHead() != build.GSettings.FProtoMask.FNetwork {
			return errors.New("network mask is invalid")
		}

		pubKey, decMsg, err := client.DecryptMessage(
			asymmetric.NewMapPubKeys(privKey.GetPubKey()),
			netMsg.GetPayload().GetBody(),
		)
		if err != nil {
			return err
		}

		if !bytes.Equal(pubKey.ToBytes(), client.GetPrivKey().GetPubKey().ToBytes()) {
			return errors.New("got invalid public key")
		}

		pld := payload.LoadPayload64(decMsg)
		if pld.GetHead() > messageCount {
			return errors.New("got invalid head value")
		}

		if string(pld.GetBody()) != "hello, world!" {
			return errors.New("got invalid body value")
		}
	}

	return nil
}

func hashIsExist(hash []byte, listHashes []string) bool {
	strHash := encoding.HexEncode(hash)
	for _, h := range listHashes {
		if strHash == h {
			return true
		}
	}
	return false
}
