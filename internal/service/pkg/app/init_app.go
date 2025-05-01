package app

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/internal/utils/build"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/internal/utils/privkey"

	"github.com/number571/hidden-lake/internal/service/pkg/app/config"
	pkg_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

// initApp work with the raw data = read files, read args
func InitApp(pArgs []string, pFlags flag.IFlags) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(pFlags.Get("-p").GetStringValue(pArgs), "/")
	if err := build.SetBuildByPath(inputPath); err != nil {
		return nil, errors.Join(ErrSetNetworks, err)
	}

	cfgPath := filepath.Join(inputPath, pkg_settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil, pFlags.Get("-n").GetStringValue(pArgs))
	if err != nil {
		return nil, errors.Join(ErrInitConfig, err)
	}

	keyPath := filepath.Join(inputPath, pkg_settings.CPathKey)
	privKey, err := privkey.GetPrivKey(keyPath)
	if err != nil {
		return nil, errors.Join(ErrGetPrivateKey, err)
	}

	return NewApp(cfg, privKey, inputPath), nil
}
