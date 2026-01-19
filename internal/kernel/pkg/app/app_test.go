package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	testutils_gopeer "github.com/number571/go-peer/test/utils"
	"github.com/number571/hidden-lake/internal/kernel/pkg/app/config"
	"github.com/number571/hidden-lake/internal/kernel/pkg/client"
	pkg_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
	testutils "github.com/number571/hidden-lake/test/utils"
)

var (
	tgFlags = flag.NewFlagsBuilder(
		flag.NewFlagBuilder("-v", "--version").
			WithDescription("print version of application"),
		flag.NewFlagBuilder("-h", "--help").
			WithDescription("print information about application"),
		flag.NewFlagBuilder("-p", "--path").
			WithDescription("set path to config, database files").
			WithDefinedValue("."),
		flag.NewFlagBuilder("-n", "--network").
			WithDescription("set network key of connections from build").
			WithDefinedValue(""),
	).Build()
)

const (
	tcTestdataPath = "./testdata/"
	tcPathDB       = pkg_settings.CPathDB
	tcPathConfig   = pkg_settings.CPathYML
	tcPathKey      = pkg_settings.CPathKey
)

func testDeleteFiles(prefixPath string) {
	_ = os.RemoveAll(prefixPath + tcPathDB)
	_ = os.RemoveAll(prefixPath + tcPathConfig)
	_ = os.RemoveAll(prefixPath + tcPathKey)
}

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestInitApp(t *testing.T) {
	t.Parallel()

	testDeleteFiles(tcTestdataPath)
	defer testDeleteFiles(tcTestdataPath)

	if _, err := InitApp([]string{"--path", tcTestdataPath}, tgFlags); err != nil {
		t.Fatal(err)
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
			FExternal: testutils.TgAddrs[2],
			FInternal: testutils.TgAddrs[3],
		},
		FFriends: map[string]string{
			"Alice": asymmetric.NewPrivKey().GetPubKey().ToString(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	privKey := asymmetric.NewPrivKey()
	app := NewApp(cfg, privKey, ".")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := app.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			t.Error(err)
		}
	}()

	client := client.NewClient(
		client.NewBuilder(),
		client.NewRequester(
			testutils.TgAddrs[3],
			&http.Client{Timeout: time.Minute},
		),
	)

	err1 := testutils_gopeer.TryN(
		50,
		10*time.Millisecond,
		func() error {
			_, err := client.GetIndex(context.Background())
			return err
		},
	)
	if err1 != nil {
		t.Fatal(err1)
	}

	// Check public key of node
	pubKey, err := client.GetPubKey(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if pubKey.ToString() != privKey.GetPubKey().ToString() {
		t.Fatalf("public keys are not equals")
	}

	// try twice running
	go func() {
		if err := app.Run(ctx); err == nil {
			t.Error("success double run")
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
		}
	}()
	time.Sleep(100 * time.Millisecond)
}
