package app

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	hlf_settings "github.com/number571/hidden-lake/internal/applications/filesharer/pkg/settings"
	hlm_settings "github.com/number571/hidden-lake/internal/applications/messenger/pkg/settings"
	"github.com/number571/hidden-lake/internal/composite/internal/config"
	hlc_settings "github.com/number571/hidden-lake/internal/composite/pkg/settings"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

const (
	tcTestdataPath  = "./testdata/"
	tcPathConfigHLC = hlc_settings.CPathYML
	tcPathConfigHLS = hls_settings.CPathYML
	tcPathConfigHLM = hlm_settings.CPathYML
	tcPathConfigHLF = hlf_settings.CPathYML
	tcPathKeyHLS    = hls_settings.CPathKey
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
	os.RemoveAll(prefixPath + tcPathConfigHLF)
	os.RemoveAll(prefixPath + tcPathConfigHLM)
	os.RemoveAll(prefixPath + tcPathKeyHLS)
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

	app := NewApp(cfg, nil)

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

	if _, err := InitApp([]string{"path", tcTestdataPath}); err != nil {
		t.Error(err)
		return
	}

	if _, err := InitApp([]string{"path", tcTestdataPath, "parallel", "abc"}); err == nil {
		t.Error("success init app with parallel=abc")
		return
	}

	if _, err := InitApp([]string{"path", "./not_exist/path/to/hle"}); err == nil {
		t.Error("success init app with undefined dir key")
		return
	}
}
