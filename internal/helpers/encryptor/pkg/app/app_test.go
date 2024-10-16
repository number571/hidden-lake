package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
	testutils "github.com/number571/go-peer/test/utils"
	"github.com/number571/hidden-lake/internal/helpers/encryptor/internal/config"
	"github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/client"
	pkg_settings "github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/settings"
)

const (
	tcTestdataPath1024 = "./testdata/1024"
	tcTestdataPath4096 = "./testdata/4096"
	tcPathConfig       = pkg_settings.CPathYML
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

func testDeleteFiles(prefixPath string) {
	os.RemoveAll(prefixPath + tcPathConfig)
}

func TestApp(t *testing.T) {
	t.Parallel()

	testDeleteFiles("./")
	defer testDeleteFiles("./")

	// Run application
	cfg, err := config.BuildConfig(tcPathConfig, &config.SConfig{
		FSettings: &config.SConfigSettings{
			FMessageSizeBytes: testutils.TCMessageSize,
			FWorkSizeBits:     testutils.TCWorkSize,
			FKeySizeBits:      testutils.TcKeySize,
			FNetworkKey:       "_",
		},
		FAddress: &config.SAddress{
			FHTTP: testutils.TgAddrs[55],
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	privKey := asymmetric.LoadRSAPrivKey(testutils.TcPrivKey1024)
	app := NewApp(cfg, privKey, 1)

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
			"http://"+testutils.TgAddrs[55],
			&http.Client{Timeout: time.Minute},
			testNetworkMessageSettings(),
		),
	)

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

func TestInitApp(t *testing.T) {
	t.Parallel()

	testDeleteFiles(tcTestdataPath1024)
	testDeleteFiles(tcTestdataPath4096)

	defer testDeleteFiles(tcTestdataPath1024)
	defer testDeleteFiles(tcTestdataPath4096)

	if _, err := InitApp([]string{}, tcTestdataPath4096, 1); err != nil {
		t.Error(err)
		return
	}

	if _, err := InitApp([]string{"parallel", "abc"}, tcTestdataPath4096, 1); err == nil {
		t.Error("success init app with parallel=abc")
		return
	}

	if _, err := InitApp([]string{}, tcTestdataPath1024, 1); err == nil {
		t.Error("success init app with diff key size")
		return
	}

	if _, err := InitApp([]string{}, "./not_exist/path/to/hle", 1); err == nil {
		t.Error("success init app with undefined dir key")
		return
	}
}

func testNetworkMessageSettings() net_message.ISettings {
	return net_message.NewSettings(&net_message.SSettings{
		FNetworkKey:   testutils.TCNetworkKey,
		FWorkSizeBits: testutils.TCWorkSize,
	})
}
