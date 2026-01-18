package app

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/pkg/types"
	build "github.com/number571/hidden-lake/build/environment"
	"github.com/number571/hidden-lake/internal/utils/flag"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	"github.com/number571/hidden-lake/internal/utils/privkey"

	"github.com/number571/hidden-lake/internal/kernel/pkg/app/config"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
)

// initApp work with the raw data = read files, read args
func InitApp(pArgs []string, pFlags flag.IFlags) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(pFlags.Get("-p").GetStringValue(pArgs), "/")
	if err := os.MkdirAll(inputPath, 0700); err != nil {
		return nil, errors.Join(ErrMkdirPath, err)
	}

	okLoaded, err := build.SetBuildByPath(inputPath)
	if err != nil {
		return nil, errors.Join(ErrSetBuild, err)
	}

	cfgPath := filepath.Join(inputPath, hlk_settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil, pFlags.Get("-n").GetStringValue(pArgs))
	if err != nil {
		return nil, errors.Join(ErrInitConfig, err)
	}

	stdfLogger := std_logger.NewStdLogger(cfg.GetLogging(), std_logger.GetLogFunc())
	build.LogLoadedBuildFiles(hlk_settings.GetAppShortNameFMT(), stdfLogger, okLoaded)

	keyPath := filepath.Join(inputPath, hlk_settings.CPathKey)
	privKey, err := privkey.GetPrivKey(keyPath)
	if err != nil {
		return nil, errors.Join(ErrGetPrivateKey, err)
	}

	return NewApp(cfg, privKey, inputPath), nil
}
