package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	testutils "github.com/number571/go-peer/test/utils"
	"github.com/number571/hidden-lake/internal/service/pkg/app/config"
	"github.com/number571/hidden-lake/internal/service/pkg/client"
	pkg_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

const (
	tcTestdataPath = "./testdata/"
	tcPathDB       = pkg_settings.CPathDB
	tcPathConfig   = pkg_settings.CPathYML
	tcPathKey      = pkg_settings.CPathKey
)

func testDeleteFiles(prefixPath string) {
	os.RemoveAll(prefixPath + tcPathDB)
	os.RemoveAll(prefixPath + tcPathConfig)
	os.RemoveAll(prefixPath + tcPathKey)
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

func TestInitApp(t *testing.T) {
	t.Parallel()

	testDeleteFiles(tcTestdataPath)
	defer testDeleteFiles(tcTestdataPath)

	if _, err := InitApp([]string{"path", tcTestdataPath}); err != nil {
		t.Error(err)
		return
	}

	if _, err := InitApp([]string{"path", tcTestdataPath, "parallel", "abc"}); err == nil {
		t.Error("success init app with parallel=abc")
		return
	}

	if _, err := InitApp([]string{"path", "./not_exist/path/to/hls"}); err == nil {
		t.Error("success init app with undefined dir key")
		return
	}
}

func TestApp(t *testing.T) {
	t.Parallel()

	testDeleteFiles("./")
	defer testDeleteFiles("./")

	// Run application
	cfg, err := config.BuildConfig(tcPathConfig, &config.SConfig{
		FSettings: &config.SConfigSettings{
			FMessageSizeBytes: (8 << 10),
			FWorkSizeBits:     10,
			FQueuePeriodMS:    5_000,
			FFetchTimeoutMS:   30_000,
			FNetworkKey:       "_",
		},
		FAddress: &config.SAddress{
			FTCP:  testutils.TgAddrs[2],
			FHTTP: testutils.TgAddrs[3],
		},
		FFriends: map[string]string{
			"Alice": asymmetric.NewPrivKey().GetPubKey().ToString(),
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	privKey := asymmetric.NewPrivKey()
	app := NewApp(cfg, privKey, ".", 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := app.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			t.Error(err)
			return
		}
	}()

	client := client.NewClient(
		client.NewBuilder(),
		client.NewRequester(
			"http://"+testutils.TgAddrs[3],
			&http.Client{Timeout: time.Minute},
		),
	)

	err1 := testutils.TryN(
		50,
		10*time.Millisecond,
		func() error {
			_, err := client.GetIndex(context.Background())
			return err
		},
	)
	if err1 != nil {
		t.Error(err1)
		return
	}

	// Check public key of node
	pubKey, err := client.GetPubKey(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	if pubKey.ToString() != privKey.GetPubKey().ToString() {
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
