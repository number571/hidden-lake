package app

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/types"
	hla_tcp_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
	"github.com/number571/hidden-lake/internal/composite/pkg/app/config"
	hlc_settings "github.com/number571/hidden-lake/internal/composite/pkg/settings"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	hls_messenger_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	hls_pinger_settings "github.com/number571/hidden-lake/internal/services/pinger/pkg/settings"
	"github.com/number571/hidden-lake/pkg/utils/flag"
)

var (
	tgFlags = flag.NewFlagsBuilder(
		flag.NewFlagBuilder("-v", "--version").
			WithDescription("print version of service"),
		flag.NewFlagBuilder("-h", "--help").
			WithDescription("print information about service"),
		flag.NewFlagBuilder("-p", "--path").
			WithDescription("set path to config, database files").
			WithDefinedValue("."),
		flag.NewFlagBuilder("-n", "--network").
			WithDescription("set network key of connections from build").
			WithDefinedValue(""),
	).Build()
)

const (
	tcTestdataPath     = "./testdata/"
	tcPathConfigHLC    = hlc_settings.CPathYML
	tcPathConfigHLS    = hlk_settings.CPathYML
	tcPathConfigHLM    = hls_messenger_settings.CPathYML
	tcPathConfigHLP    = hls_pinger_settings.CPathYML
	tcPathConfigHLAtcp = hla_tcp_settings.CPathYML
	tcPathKeyHLS       = hlk_settings.CPathKey
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
	_ = os.RemoveAll(prefixPath + tcPathConfigHLC)
	_ = os.RemoveAll(prefixPath + tcPathConfigHLS)
	_ = os.RemoveAll(prefixPath + tcPathConfigHLP)
	_ = os.RemoveAll(prefixPath + tcPathConfigHLM)
	_ = os.RemoveAll(prefixPath + tcPathConfigHLAtcp)
	_ = os.RemoveAll(prefixPath + tcPathKeyHLS)
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
		t.Fatal(err)
	}

	app := NewApp(cfg, []types.IRunner{&tsRunner{}})

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

	if _, err := InitApp([]string{"--path", "./not_exist/path/to/hlc"}, tgFlags); err == nil {
		t.Fatal("success init app with undefined dir key")
	}
}
