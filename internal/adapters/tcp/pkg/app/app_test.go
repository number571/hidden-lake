package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/payload"
	testutils_gopeer "github.com/number571/go-peer/test/utils"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/adapters/tcp/pkg/app/config"
	"github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/pkg/adapters/http/client"
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
		flag.NewFlagBuilder("-n", "--network").
			WithDescription("set network key of connections from build").
			WithDefinedValue(""),
	).Build()
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SAppError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

const tcPathConfig = "./testdata/"
const tcDataConfig = `settings:
  message_size_bytes: 8192
`

func TestInitApp(t *testing.T) {
	t.Parallel()

	testDeleteFiles(tcPathConfig)
	defer testDeleteFiles(tcPathConfig)

	if err := os.WriteFile(tcPathConfig+"hla_tcp.yml", []byte(tcDataConfig), 0600); err != nil {
		t.Fatal(err)
	}

	app, err := InitApp([]string{"--path", tcPathConfig}, tgFlags)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := app.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			t.Error(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()
}

func testDeleteFiles(path string) {
	_ = os.RemoveAll(path + settings.CPathYML)
	_ = os.RemoveAll(path + settings.CPathDB)
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
			FExternal: testutils.TgAddrs[20],
		},
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
		t.Fatal(err1)
	}

	msgBytes := []byte("hello, world!")
	msgBytes = append(msgBytes, random.NewRandom().GetBytes(uint64(8192-len(msgBytes)))...) //nolint:gosec
	netMsg := layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{
				FWorkSizeBits: 10,
				FNetworkKey:   "_",
			}),
		}),
		payload.NewPayload32(build.GetSettings().FProtoMask.FNetwork, msgBytes),
	)

	if err := client.ProduceMessage(ctx, netMsg); err != nil {
		t.Fatal(err)
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
