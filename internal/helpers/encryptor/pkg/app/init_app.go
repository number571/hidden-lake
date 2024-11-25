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
func InitApp(pArgs []string, pFlags flag.IFlags) (types.IRunner, error) {
	strParallel := pFlags.Get("threads").GetStringValue(pArgs)
	setParallel, err := strconv.ParseUint(strParallel, 10, 64)
	if err != nil {
		return nil, errors.Join(ErrGetParallelValue, err)
	}

	inputPath := strings.TrimSuffix(pFlags.Get("path").GetStringValue(pArgs), "/")

	cfgPath := filepath.Join(inputPath, settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil, pFlags.Get("network").GetStringValue(pArgs))
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
