package app

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/types"
	hla_tcp_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
	hlm_settings "github.com/number571/hidden-lake/internal/applications/messenger/pkg/settings"
	hlp_settings "github.com/number571/hidden-lake/internal/applications/pinger/pkg/settings"
	"github.com/number571/hidden-lake/internal/composite/pkg/app/config"
	hlc_settings "github.com/number571/hidden-lake/internal/composite/pkg/settings"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
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
		flag.NewFlagBuilder("-n", "--network").
			WithDescription("set network key for connections").
			WithDefinedValue(""),
		flag.NewFlagBuilder("-t", "--threads").
			WithDescription("set num of parallel functions to calculate PoW").
			WithDefinedValue("1"),
	).Build()
)

const (
	tcTestdataPath     = "./testdata/"
	tcPathConfigHLC    = hlc_settings.CPathYML
	tcPathConfigHLS    = hls_settings.CPathYML
	tcPathConfigHLM    = hlm_settings.CPathYML
	tcPathConfigHLP    = hlp_settings.CPathYML
	tcPathConfigHLAtcp = hla_tcp_settings.CPathYML
	tcPathKeyHLS       = hls_settings.CPathKey
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
	os.RemoveAll(prefixPath + tcPathConfigHLC)
	os.RemoveAll(prefixPath + tcPathConfigHLS)
	os.RemoveAll(prefixPath + tcPathConfigHLP)
	os.RemoveAll(prefixPath + tcPathConfigHLM)
	os.RemoveAll(prefixPath + tcPathConfigHLAtcp)
	os.RemoveAll(prefixPath + tcPathKeyHLS)
}

type tsRunner struct{}

func (p *tsRunner) Run(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func TestApp(t *testing.T) {
	t.Parallel()

	testDeleteFiles("./")
	defer testDeleteFiles("./")

	// Run application
	cfg, err := config.BuildConfig(tcPathConfigHLC, &config.SConfig{
		FServices: []string{"test"},
	})
	if err != nil {
		t.Error(err)
		return
	}

	app := NewApp(cfg, []types.IRunner{&tsRunner{}})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := app.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			t.Error(err)
			return
		}
	}()

	time.Sleep(100 * time.Millisecond)

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

	if _, err := InitApp([]string{"--path", tcTestdataPath}, tgFlags); err != nil {
		t.Error(err)
		return
	}

	if _, err := InitApp([]string{"--path", tcTestdataPath, "--threads", "abc"}, tgFlags); err == nil {
		t.Error("success init app with threads=abc")
		return
	}

	if _, err := InitApp([]string{"--path", "./not_exist/path/to/hlc"}, tgFlags); err == nil {
		t.Error("success init app with undefined dir key")
		return
	}
}
