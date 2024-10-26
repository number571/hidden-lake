package app

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/go-peer/pkg/utils"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/internal/utils/privkey"

	"github.com/number571/hidden-lake/internal/service/internal/config"
	pkg_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

// initApp work with the raw data = read files, read args
func InitApp(pArgs []string) (types.IRunner, error) {
	strParallel := flag.GetFlagValue(pArgs, "parallel", "1")
	setParallel, err := strconv.ParseUint(strParallel, 10, 64)
	if err != nil {
		return nil, utils.MergeErrors(ErrGetParallel, err)
	}

	inputPath := strings.TrimSuffix(flag.GetFlagValue(pArgs, "path", "."), "/")

	cfgPath := filepath.Join(inputPath, pkg_settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil, flag.GetFlagValue(pArgs, "network", ""))
	if err != nil {
		return nil, utils.MergeErrors(ErrInitConfig, err)
	}

	keyPath := filepath.Join(inputPath, pkg_settings.CPathKey)
	privKey, err := privkey.GetPrivKey(keyPath)
	if err != nil {
		return nil, utils.MergeErrors(ErrGetPrivateKey, err)
	}

	return NewApp(cfg, privKey, inputPath, setParallel), nil
}
