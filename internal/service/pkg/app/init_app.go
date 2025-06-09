package app

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/internal/utils/build"
	"github.com/number571/hidden-lake/internal/utils/flag"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	"github.com/number571/hidden-lake/internal/utils/privkey"

	"github.com/number571/hidden-lake/internal/service/pkg/app/config"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

// initApp work with the raw data = read files, read args
func InitApp(pArgs []string, pFlags flag.IFlags) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(pFlags.Get("-p").GetStringValue(pArgs), "/")
	okLoaded, err := build.SetBuildByPath(inputPath)
	if err != nil {
		return nil, errors.Join(ErrSetNetworks, err)
	}

	cfgPath := filepath.Join(inputPath, hls_settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil, pFlags.Get("-n").GetStringValue(pArgs))
	if err != nil {
		return nil, errors.Join(ErrInitConfig, err)
	}

	stdfLogger := std_logger.NewStdLogger(cfg.GetLogging(), std_logger.GetLogFunc())
	build.LogLoadedBuildFiles(hls_settings.GetServiceName(), stdfLogger, okLoaded)

	keyPath := filepath.Join(inputPath, hls_settings.CPathKey)
	privKey, err := privkey.GetPrivKey(keyPath)
	if err != nil {
		return nil, errors.Join(ErrGetPrivateKey, err)
	}

	return NewApp(cfg, privKey, inputPath), nil
}
