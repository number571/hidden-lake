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
	"github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/app/config"
	"github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/client"
	pkg_settings "github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
	testutils "github.com/number571/hidden-lake/test/utils"
)

var (
	tgFlags = flag.NewFlagsBuilder(
		flag.NewFlagBuilder("v", "version").
			WithDescription("print information about service"),
		flag.NewFlagBuilder("h", "help").
			WithDescription("print version of service"),
		flag.NewFlagBuilder("p", "path").
			WithDescription("set path to config, database files").
			WithDefaultValue("."),
		flag.NewFlagBuilder("n", "network").
			WithDescription("set network key for connections").
			WithDefaultValue(""),
		flag.NewFlagBuilder("t", "threads").
			WithDescription("set num of parallel functions to calculate PoW").
			WithDefaultValue("1"),
	).Build()
)

const (
	tcTestdataPath = "./testdata/"
	tcPathConfig   = pkg_settings.CPathYML
	tcPathKey      = pkg_settings.CPathKey
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
	os.RemoveAll(prefixPath + tcPathKey)
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
			FNetworkKey:       "_",
		},
		FAddress: &config.SAddress{
			FInternal: testutils.TgAddrs[30],
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	privKey := asymmetric.NewPrivKey()
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
		client.NewBuilder(),
		client.NewRequester(
			testutils.TgAddrs[30],
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

	testDeleteFiles(tcTestdataPath)
	defer testDeleteFiles(tcTestdataPath)

	if _, err := InitApp([]string{"path", tcTestdataPath}, tgFlags); err != nil {
		t.Error(err)
		return
	}

	if _, err := InitApp([]string{"path", tcTestdataPath, "threads", "abc"}, tgFlags); err == nil {
		t.Error("success init app with threads=abc")
		return
	}

	if _, err := InitApp([]string{"path", "./not_exist/path/to/hle"}, tgFlags); err == nil {
		t.Error("success init app with undefined dir key")
		return
	}
}

func testNetworkMessageSettings() net_message.ISettings {
	return net_message.NewSettings(&net_message.SSettings{
		FNetworkKey:   "_",
		FWorkSizeBits: 10,
	})
}
