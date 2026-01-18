package app

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/pkg/types"
	build "github.com/number571/hidden-lake/build/environment"
	"github.com/number571/hidden-lake/internal/services/pinger/pkg/app/config"
	hls_pinger_settings "github.com/number571/hidden-lake/internal/services/pinger/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func InitApp(pArgs []string, pFlags flag.IFlags) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(pFlags.Get("-p").GetStringValue(pArgs), "/")
	if err := os.MkdirAll(inputPath, 0700); err != nil {
		return nil, errors.Join(ErrMkdirPath, err)
	}

	okLoaded, err := build.SetBuildByPath(inputPath)
	if err != nil {
		return nil, errors.Join(ErrSetBuild, err)
	}

	cfgPath := filepath.Join(inputPath, hls_pinger_settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil)
	if err != nil {
		return nil, errors.Join(ErrInitConfig, err)
	}

	stdfLogger := std_logger.NewStdLogger(cfg.GetLogging(), std_logger.GetLogFunc())
	build.LogLoadedBuildFiles(hls_pinger_settings.GetAppShortNameFMT(), stdfLogger, okLoaded)

	return NewApp(cfg), nil
}
