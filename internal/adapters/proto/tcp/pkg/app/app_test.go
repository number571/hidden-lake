package app

// import (
// 	"context"
// 	"errors"
// 	"os"
// 	"testing"
// 	"time"

// 	consumer_app "github.com/number571/hidden-lake/internal/adapters/tcp/internal/consumer/pkg/app"
// 	producer_app "github.com/number571/hidden-lake/internal/adapters/tcp/internal/producer/pkg/app"
// 	"github.com/number571/hidden-lake/internal/adapters/tcp/pkg/app/config"
// 	"github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
// 	"github.com/number571/hidden-lake/internal/utils/flag"
// 	testutils "github.com/number571/hidden-lake/test/utils"
// )

// var (
// 	tgFlags = flag.NewFlagsBuilder(
// 		flag.NewFlagBuilder("v", "version").
// 			WithDescription("print information about service"),
// 		flag.NewFlagBuilder("h", "help").
// 			WithDescription("print version of service"),
// 		flag.NewFlagBuilder("p", "path").
// 			WithDescription("set path to config, database files").
// 			WithDefaultValue("."),
// 		flag.NewFlagBuilder("n", "network").
// 			WithDescription("set network key for connections").
// 			WithDefaultValue(""),
// 	).Build()
// )

// const (
// 	tcPathConfig = "./testdata/"
// )

// func TestError(t *testing.T) {
// 	t.Parallel()

// 	str := "value"
// 	err := &SAppError{str}
// 	if err.Error() != errPrefix+str {
// 		t.Error("incorrect err.Error()")
// 		return
// 	}
// }

// func TestInitApp(t *testing.T) {
// 	t.Parallel()

// 	testDeleteFiles(tcPathConfig)
// 	defer testDeleteFiles(tcPathConfig)

// 	if _, err := InitApp([]string{"path", tcPathConfig}, tgFlags); err != nil {
// 		t.Error(err)
// 		return
// 	}
// }

// func testDeleteFiles(path string) {
// 	os.RemoveAll(path + settings.CPathYML)
// }

// func TestApp(t *testing.T) {
// 	t.Parallel()

// 	testDeleteFiles("./")
// 	defer testDeleteFiles("./")

// 	// Run application
// 	cfg, err := config.BuildConfig(settings.CPathYML, &config.SConfig{
// 		FSettings: &config.SConfigSettings{
// 			FWorkSizeBits: 10,
// 			FNetworkKey:   "_",
// 		},
// 		FAddress: testutils.TgAddrs[45],
// 		FConnection: &config.SConnection{
// 			FHLTHost: "hlt",
// 			FSrvHost: "srv",
// 		},
// 	})
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	app := NewApp(
// 		cfg,
// 		consumer_app.NewApp(cfg),
// 		producer_app.NewApp(cfg),
// 	)

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	go func() {
// 		if err := app.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
// 			t.Error(err)
// 			return
// 		}
// 	}()

// 	time.Sleep(100 * time.Millisecond)

// 	// try twice running
// 	go func() {
// 		if err := app.Run(ctx); err == nil {
// 			t.Error("success double run")
// 			return
// 		}
// 	}()

// 	time.Sleep(100 * time.Millisecond)
// 	cancel()
// 	time.Sleep(100 * time.Millisecond)

// 	ctx1, cancel1 := context.WithCancel(context.Background())
// 	defer cancel1()

// 	// try twice running
// 	go func() {
// 		if err := app.Run(ctx1); err != nil && !errors.Is(err, context.Canceled) {
// 			t.Error(err)
// 			return
// 		}
// 	}()
// 	time.Sleep(100 * time.Millisecond)
// }
