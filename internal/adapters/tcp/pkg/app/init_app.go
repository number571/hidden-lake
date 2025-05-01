package app

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/internal/adapters/tcp/pkg/app/config"
	"github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/build"
	"github.com/number571/hidden-lake/internal/utils/flag"
)

func InitApp(pArgs []string, pFlags flag.IFlags) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(pFlags.Get("-p").GetStringValue(pArgs), "/")
	if err := build.SetBuildByPath(inputPath); err != nil {
		return nil, errors.Join(ErrSetNetworks, err)
	}

	cfgPath := filepath.Join(inputPath, settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil, pFlags.Get("-n").GetStringValue(pArgs))
	if err != nil {
		return nil, errors.Join(ErrInitConfig, err)
	}

	return NewApp(cfg, inputPath), nil
}
