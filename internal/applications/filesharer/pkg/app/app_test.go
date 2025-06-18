package app

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/number571/hidden-lake/internal/applications/filesharer/pkg/app/config"
	pkg_settings "github.com/number571/hidden-lake/internal/applications/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
	testutils "github.com/number571/hidden-lake/test/utils"
)

var (
	tgFlags = flag.NewFlagsBuilder(
		flag.NewFlagBuilder("-v", "--version").
			WithDescription("print information about service"),
		flag.NewFlagBuilder("-h", "--help").
			WithDescription("print version of service"),
		flag.NewFlagBuilder("-p", "--path").
			WithDescription("set path to config, database files").
			WithDefinedValue("."),
	).Build()
)

const (
	tcTestdataPath = "./testdata/"
	tcPathConfig   = pkg_settings.CPathYML
	tcPathStorage  = pkg_settings.CPathSTG
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SAppError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func testDeleteFiles(prefixPath string) {
	_ = os.RemoveAll(prefixPath + tcPathStorage)
	_ = os.RemoveAll(prefixPath + tcPathConfig)
}

func TestApp(t *testing.T) {
	t.Parallel()

	testDeleteFiles("./")
	defer testDeleteFiles("./")

	// Run application
	cfg, err := config.BuildConfig(tcPathConfig, &config.SConfig{
		FSettings: &config.SConfigSettings{
			FPageOffset: 10,
		},
		FAddress: &config.SAddress{
			FInternal: testutils.TgAddrs[37],
			FExternal: testutils.TgAddrs[29],
		},
		FConnection: "test_connection",
	})
	if err != nil {
		t.Fatal(err)
	}

	app := NewApp(cfg, ".")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := app.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			t.Error(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

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

func TestInitApp(t *testing.T) {
	t.Parallel()

	testDeleteFiles(tcTestdataPath)
	defer testDeleteFiles(tcTestdataPath)

	if _, err := InitApp([]string{"--path", tcTestdataPath}, tgFlags); err != nil {
		t.Fatal(err)
	}
}
