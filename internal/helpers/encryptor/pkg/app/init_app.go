package app

import (
	"errors"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/app/config"
	"github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/internal/utils/privkey"
)

// initApp work with the raw data = read files, read args
func InitApp(pArgs []string) (types.IRunner, error) {
	strParallel := flag.GetFlagValue(pArgs, "parallel", "1")
	setParallel, err := strconv.ParseUint(strParallel, 10, 64)
	if err != nil {
		return nil, errors.Join(ErrGetParallelValue, err)
	}

	inputPath := strings.TrimSuffix(flag.GetFlagValue(pArgs, "path", "."), "/")

	cfgPath := filepath.Join(inputPath, settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil, flag.GetFlagValue(pArgs, "network", ""))
	if err != nil {
		return nil, errors.Join(ErrInitConfig, err)
	}

	keyPath := filepath.Join(inputPath, settings.CPathKey)
	privKey, err := privkey.GetPrivKey(keyPath)
	if err != nil {
		return nil, errors.Join(ErrGetPrivateKey, err)
	}

	return NewApp(cfg, privKey, setParallel), nil
}
