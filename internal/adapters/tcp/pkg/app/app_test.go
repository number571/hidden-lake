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
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/storage/cache"
	testutils_gopeer "github.com/number571/go-peer/test/utils"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/adapters/tcp/pkg/app/config"
	"github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/pkg/api/adapters/http/client"
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

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SError{str}
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

	if err := os.WriteFile(tcPathConfig+"hla-tcp.yml", []byte(tcDataConfig), 0600); err != nil {
		t.Fatal(err)
	}

	app, err := InitApp([]string{"--path", tcPathConfig}, tgFlags)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := InitApp([]string{"--path", tcPathConfig + "/failed"}, tgFlags); err == nil {
		t.Fatal("success init app with invalid config")
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
			err := client.GetIndex(context.Background(), settings.CAppAdapterName)
			return err
		},
	)
	if err1 != nil {
		t.Fatal(err1)
	}

	layer1Settings := layer1.NewSettings(&layer1.SSettings{
		FWorkSizeBits: 10,
		FNetworkKey:   "_",
	})

	msgBytes := []byte("hello, world!")
	msgBytes = append(msgBytes, random.NewRandom().GetBytes(uint64(8192-len(msgBytes)))...) //nolint:gosec
	netMsg := layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1Settings,
		}),
		payload.NewPayload32(build.GetSettings().FProtoMask.FNetwork, msgBytes),
	)

	if err := client.ProduceMessage(ctx, netMsg); err != nil {
		t.Fatal(err)
	}

	netNode := network.NewNode(
		network.NewSettings(&network.SSettings{
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FLimitMessageSizeBytes: 8192,
				FWaitReadTimeout:       time.Second,
				FDialTimeout:           time.Second,
				FReadTimeout:           time.Second,
				FWriteTimeout:          time.Second,
				FMessageSettings:       layer1Settings,
			}),
			FAddress:      "",
			FMaxConnects:  1,
			FReadTimeout:  time.Second,
			FWriteTimeout: time.Second,
		}),
		cache.NewLRUCache(128),
	)

	if err := netNode.AddConnection(ctx, testutils.TgAddrs[20]); err != nil {
		t.Fatal(err)
	}

	if err := netNode.BroadcastMessage(ctx, netMsg); err != nil {
		t.Fatal(err)
	}

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
