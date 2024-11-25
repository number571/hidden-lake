package app

import (
	"errors"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/internal/utils/privkey"

	"github.com/number571/hidden-lake/internal/service/pkg/app/config"
	pkg_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

// initApp work with the raw data = read files, read args
func InitApp(pArgs []string) (types.IRunner, error) {
	strParallel := flag.GetStringFlagValue(pArgs, []string{"t", "threads"}, "1")
	setParallel, err := strconv.ParseUint(strParallel, 10, 64)
	if err != nil {
		return nil, errors.Join(ErrGetParallel, err)
	}

	inputPath := strings.TrimSuffix(flag.GetStringFlagValue(pArgs, []string{"p", "path"}, "."), "/")

	cfgPath := filepath.Join(inputPath, pkg_settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil, flag.GetStringFlagValue(pArgs, []string{"n", "network"}, ""))
	if err != nil {
		return nil, errors.Join(ErrInitConfig, err)
	}

	keyPath := filepath.Join(inputPath, pkg_settings.CPathKey)
	privKey, err := privkey.GetPrivKey(keyPath)
	if err != nil {
		return nil, errors.Join(ErrGetPrivateKey, err)
	}

	return NewApp(cfg, privKey, inputPath, setParallel), nil
}
