package app

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/go-peer/pkg/utils"
	"github.com/number571/hidden-lake/internal/helpers/encryptor/internal/config"
	"github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/internal/utils/privkey"
)

// initApp work with the raw data = read files, read args
func InitApp(pArgs []string, pDefaultPath string, pDefaultParallel uint64) (types.IRunner, error) {
	strParallel := flag.GetFlagValue(pArgs, "parallel", strconv.FormatUint(pDefaultParallel, 10))
	setParallel, err := strconv.ParseUint(strParallel, 10, 64)
	if err != nil {
		return nil, utils.MergeErrors(ErrGetParallelValue, err)
	}

	inputPath := strings.TrimSuffix(flag.GetFlagValue(pArgs, "path", pDefaultPath), "/")

	cfgPath := filepath.Join(inputPath, settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil)
	if err != nil {
		return nil, utils.MergeErrors(ErrInitConfig, err)
	}

	keyPath := filepath.Join(inputPath, settings.CPathKey)
	privKey, err := privkey.GetPrivKey(keyPath)
	if err != nil {
		return nil, utils.MergeErrors(ErrGetPrivateKey, err)
	}

	return NewApp(cfg, privKey, setParallel), nil
}
