package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	testutils "github.com/number571/go-peer/test/utils"
	"github.com/number571/hidden-lake/internal/helpers/loader/internal/config"
	"github.com/number571/hidden-lake/internal/helpers/loader/pkg/client"
	pkg_settings "github.com/number571/hidden-lake/internal/helpers/loader/pkg/settings"
)

const (
	tcPathConfig = pkg_settings.CPathYML
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SAppError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func testDeleteFiles() {
	os.RemoveAll(tcPathConfig)
}

func TestApp(t *testing.T) {
	t.Parallel()

	testDeleteFiles()
	defer testDeleteFiles()

	// Run application
	cfg, err := config.BuildConfig(tcPathConfig, &config.SConfig{
		FSettings: &config.SConfigSettings{
			FMessagesCapacity: testutils.TCCapacity,
			FWorkSizeBits:     testutils.TCWorkSize,
			FNetworkKey:       "_",
		},
		FAddress: &config.SAddress{
			FHTTP: testutils.TgAddrs[56],
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	app := NewApp(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := app.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			t.Error(err)
			return
		}
	}()

	time.Sleep(100 * time.Millisecond)
	client := client.NewClient(
		client.NewRequester(
			"http://"+testutils.TgAddrs[56],
			&http.Client{Timeout: time.Minute},
		),
	)

	// Check public key of node
	index, err := client.GetIndex(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	if index != pkg_settings.CServiceFullName {
		t.Errorf("public keys are not equals")
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
