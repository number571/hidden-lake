package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/hidden-lake/build"
	hll_client "github.com/number571/hidden-lake/internal/helpers/loader/pkg/client"
	hls_app "github.com/number571/hidden-lake/internal/helpers/traffic/pkg/app"
	hlt_client "github.com/number571/hidden-lake/internal/helpers/traffic/pkg/client"
	"github.com/number571/hidden-lake/internal/utils/flag"
)

var (
	tgFlags = flag.NewFlags(
		flag.NewFlagBuilder("v", "version").
			WithDescription("print information about service").
			Build(),
		flag.NewFlagBuilder("h", "help").
			WithDescription("print version of service").
			Build(),
		flag.NewFlagBuilder("p", "path").
			WithDescription("set path to config, database files").
			WithDefaultValue(".").
			Build(),
		flag.NewFlagBuilder("n", "network").
			WithDescription("set network key for connections").
			WithDefaultValue("").
			Build(),
	)
)

const (
	tcTestData = "./test_data"
	tcNameHLT1 = tcTestData + "/hlt_1"
	tcNameHLT2 = tcTestData + "/hlt_2"
)

func testCreateHLT(
	netMsgSettings net_message.ISettings,
	path string,
	addr string,
) (context.CancelFunc, hlt_client.IClient, error) {
	ctx, cancel := context.WithCancel(context.Background())

	if err := copyWithPaste(path, addr); err != nil {
		return cancel, nil, err
	}

	app1, err := hls_app.InitApp([]string{"path", path}, tgFlags)
	if err != nil {
		return cancel, nil, err
	}

	go func() { _ = app1.Run(ctx) }()
	time.Sleep(100 * time.Millisecond)

	hltClient1 := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			"http://"+addr,
			&http.Client{Timeout: time.Minute / 2},
			netMsgSettings,
		),
	)

	return cancel, hltClient1, nil
}

func testInitTransfer() {
	os.RemoveAll(tcNameHLT1)
	os.RemoveAll(tcNameHLT2)

	_ = os.Mkdir(tcNameHLT1, 0o777)
	_ = os.Mkdir(tcNameHLT2, 0o777)
}

func TestHandleTransferAPI(t *testing.T) {
	t.Parallel()

	testInitTransfer()
	defer func() {
		os.RemoveAll(tcNameHLT1)
		os.RemoveAll(tcNameHLT2)
	}()

	// INIT SERVICES

	netMsgSettings := net_message.NewConstructSettings(&net_message.SConstructSettings{
		FSettings: net_message.NewSettings(&net_message.SSettings{
			FWorkSizeBits: tcWorkSize,
			FNetworkKey:   tcNetworkKey,
		}),
	})

	cancel1, hltClient1, err := testCreateHLT(netMsgSettings.GetSettings(), tcNameHLT1, tgProducer)
	if err != nil {
		t.Error(err)
		return
	}
	defer cancel1()

	cancel2, hltClient2, err := testCreateHLT(netMsgSettings.GetSettings(), tcNameHLT2, tgConsumer)
	if err != nil {
		t.Error(err)
		return
	}
	defer cancel2()

	service := testRunService(tgTService)
	defer service.Close()

	time.Sleep(100 * time.Millisecond)
	hllClient := hll_client.NewClient(
		hll_client.NewRequester(
			"http://"+tgTService,
			&http.Client{Timeout: time.Second / 2},
		),
	)

	// PUSH MESSAGES

	privKey := asymmetric.NewPrivKey()
	client := client.NewClient(privKey, tcMessageSize)

	for i := 0; i < 5; i++ {
		encMsg, err := client.EncryptMessage(
			privKey.GetPubKey(),
			payload.NewPayload64(
				uint64(i),
				[]byte("hello, world!"),
			).ToBytes(),
		)
		if err != nil {
			t.Error(err)
			return
		}
		netMsg := net_message.NewMessage(
			netMsgSettings,
			payload.NewPayload32(
				build.GSettings.FProtoMask.FNetwork,
				encMsg,
			),
		)
		err = hltClient1.PutMessage(context.Background(), netMsg)
		if err != nil {
			t.Error(err)
			return
		}
	}

	// TRANSFER MESSAGES

	if err := hllClient.RunTransfer(context.Background()); err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Second)

	// LOAD MESSAGES

	for i := uint64(0); ; i++ {
		h, err := hltClient2.GetHash(context.Background(), i)
		if err != nil {
			if i != 5 {
				t.Error(i, err)
			}
			return
		}

		netMsg, err := hltClient2.GetMessage(context.Background(), h)
		if err != nil {
			t.Error(err)
			return
		}

		pubKey, decMsg, err := client.DecryptMessage(
			asymmetric.NewMapPubKeys(privKey.GetPubKey()),
			netMsg.GetPayload().GetBody(),
		)
		if err != nil {
			t.Error(err)
			return
		}

		pld := payload.LoadPayload64(decMsg)
		if pld.GetHead() != i {
			t.Error("got bad index")
			return
		}

		if !bytes.Equal(pubKey.ToBytes(), client.GetPrivKey().GetPubKey().ToBytes()) {
			t.Error("got bad public key")
			return
		}
	}
}

func copyWithPaste(pathTo, addr string) error {
	cfgDataFmt, err := os.ReadFile(tcTestData + "/hlt_copy.yml")
	if err != nil {
		return err
	}
	return os.WriteFile(
		pathTo+"/hlt.yml",
		[]byte(fmt.Sprintf(string(cfgDataFmt), tcWorkSize, tcNetworkKey, addr)),
		0o600,
	)
}
