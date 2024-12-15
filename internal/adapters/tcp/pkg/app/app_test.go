package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	testutils_gopeer "github.com/number571/go-peer/test/utils"
	"github.com/number571/hidden-lake/internal/adapters/tcp/pkg/app/config"
	"github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/pkg/adapters/http/client"
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
	).Build()
)

const (
	tcPathConfig = "./testdata/"
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

func TestInitApp(t *testing.T) {
	t.Parallel()

	testDeleteFiles(tcPathConfig)
	defer testDeleteFiles(tcPathConfig)

	if _, err := InitApp([]string{"path", tcPathConfig}, tgFlags); err != nil {
		t.Error(err)
		return
	}
}

func testDeleteFiles(path string) {
	os.RemoveAll(path + settings.CPathYML)
	os.RemoveAll(path + settings.CPathDB)
}

func TestApp(t *testing.T) {
	t.Parallel()

	testDeleteFiles("./")
	defer testDeleteFiles("./")

	// Run application
	cfg, err := config.BuildConfig(settings.CPathYML, &config.SConfig{
		FSettings: &config.SConfigSettings{
			FMessageSizeBytes: 8192,
			FWorkSizeBits:     10,
			FNetworkKey:       "_",
			FDatabaseEnabled:  true,
		},
		FAddress: &config.SAddress{
			FInternal: testutils.TgAddrs[17],
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	app := NewApp(cfg, ".")

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
			testutils.TgAddrs[17],
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
		t.Error(err1)
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
