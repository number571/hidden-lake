package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	hle_client "github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/client"
	"github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/settings"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestErrorsAPI(t *testing.T) {
	t.Parallel()

	client := hle_client.NewClient(
		hle_client.NewBuilder(),
		hle_client.NewRequester(
			"http://"+testutils.TcUnknownHost,
			&http.Client{Timeout: time.Second},
			testNetworkMessageSettings(),
		),
	)

	if _, err := client.EncryptMessage(context.Background(), "_", payload.NewPayload64(1, []byte{123})); err == nil {
		t.Error("success encrypt message with unknown host")
		return
	}

	pld := payload.NewPayload32(tcHead, []byte(tcBody))
	sett := message.NewConstructSettings(&message.SConstructSettings{
		FSettings: message.NewSettings(&message.SSettings{
			FWorkSizeBits: tcWorkSize,
		}),
	})
	if _, _, err := client.DecryptMessage(context.Background(), message.NewMessage(sett, pld)); err == nil {
		t.Error("success decrypt message with unknown host")
		return
	}

	if _, err := client.GetIndex(context.Background()); err == nil {
		t.Error("success get index with unknown host")
		return
	}

	if _, err := client.GetSettings(context.Background()); err == nil {
		t.Error("success get settings with unknown host")
		return
	}

	if _, err := client.GetPubKey(context.Background()); err == nil {
		t.Error("success get pub key with unknown host")
		return
	}
}

func TestHandleIndexAPI(t *testing.T) {
	t.Parallel()

	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 4)
	defer os.Remove(pathCfg)

	_, service := testRunService(pathCfg, testutils.TgAddrs[32])
	defer service.Close()

	time.Sleep(100 * time.Millisecond)
	hleClient := hle_client.NewClient(
		hle_client.NewBuilder(),
		hle_client.NewRequester(
			"http://"+testutils.TgAddrs[32],
			&http.Client{Timeout: time.Second / 2},
			testNetworkMessageSettings(),
		),
	)

	title, err := hleClient.GetIndex(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if title != settings.CServiceFullName {
		t.Error("incorrect title pattern")
		return
	}
}
