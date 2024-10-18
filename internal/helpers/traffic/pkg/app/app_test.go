package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/hidden-lake/internal/helpers/traffic/internal/config"

	hlt_client "github.com/number571/hidden-lake/internal/helpers/traffic/pkg/client"
	hlt_settings "github.com/number571/hidden-lake/internal/helpers/traffic/pkg/settings"
	testutils "github.com/number571/hidden-lake/test/utils"
)

const (
	tcPathDB     = hlt_settings.CPathDB
	tcPathConfig = hlt_settings.CPathYML
)

func testDeleteFiles() {
	os.RemoveAll(tcPathDB)
	os.RemoveAll(tcPathConfig)
}

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SAppError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestApp(t *testing.T) {
	t.Parallel()

	testDeleteFiles()
	defer testDeleteFiles()

	cfg, err := config.BuildConfig(
		tcPathConfig,
		&config.SConfig{
			FSettings: &config.SConfigSettings{
				FMessageSizeBytes: (10 << 10),
				FWorkSizeBits:     10,
				FMessagesCapacity: 32,
				FNetworkKey:       "_",
			},
			FAddress: &config.SAddress{
				FHTTP: testutils.TgAddrs[17],
			},
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := NewApp(cfg, ".")

	go func() {
		if err := app.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			t.Error(err)
			return
		}
	}()
	time.Sleep(100 * time.Millisecond)

	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			"http://"+testutils.TgAddrs[17],
			&http.Client{Timeout: time.Minute},
			net_message.NewSettings(&net_message.SSettings{
				FNetworkKey:   "_",
				FWorkSizeBits: 10,
			}),
		),
	)

	title, err := hltClient.GetIndex(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if title != hlt_settings.CServiceFullName {
		t.Error("title is incorrect")
		return
	}

	// try twice running
	go func() {
		if err := app.Run(ctx); err == nil {
			t.Error("success double run")
			return
		}
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()
	time.Sleep(100 * time.Millisecond)

	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()

	// try twice running
	go func() {
		if err := app.Run(ctx1); err != nil && !errors.Is(err, context.Canceled) {
			t.Error(err)
			return
		}
	}()
	time.Sleep(100 * time.Millisecond)
}
