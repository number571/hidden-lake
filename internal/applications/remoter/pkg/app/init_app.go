package app

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/internal/applications/remoter/pkg/app/config"
	hlr_settings "github.com/number571/hidden-lake/internal/applications/remoter/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/build"
	"github.com/number571/hidden-lake/internal/utils/flag"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func InitApp(pArgs []string, pFlags flag.IFlags) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(pFlags.Get("-p").GetStringValue(pArgs), "/")
	okLoaded, err := build.SetBuildByPath(inputPath)
	if err != nil {
		return nil, errors.Join(ErrSetNetworks, err)
	}

	cfgPath := filepath.Join(inputPath, hlr_settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil)
	if err != nil {
		return nil, errors.Join(ErrInitConfig, err)
	}

	stdfLogger := std_logger.NewStdLogger(cfg.GetLogging(), std_logger.GetLogFunc())
	build.LogLoadedBuildFiles(hlr_settings.GetServiceName(), stdfLogger, okLoaded)

	return NewApp(cfg), nil
}
