package app

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2/hybrid"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2/symmetric"
	"github.com/number571/go-peer/pkg/types"
	build "github.com/number571/hidden-lake/build/environment"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/internal/utils/keys"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	schemet "github.com/number571/hidden-lake/pkg/api/kernel/client/scheme"
	"github.com/number571/hidden-lake/pkg/network/adapters"

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

	// init default value for message size bytes (if cfg.GetSettings().GetMessageSizeBytes() == 0)
	adapterSettings := adapters.NewSettings(&adapters.SSettings{
		FMessageSizeBytes: cfg.GetSettings().GetMessageSizeBytes(),
	})

	var scheme layer2.IScheme
	msgSize := adapterSettings.GetMessageSizeBytes()

	switch cfg.GetSettings().GetCryptoSchemeType() {
	case schemet.CHybridScheme:
		scheme, err = getHybridScheme(inputPath, msgSize)
	case schemet.CSymmetricScheme:
		scheme, err = getSymmetricScheme(msgSize)
	}
	if err != nil {
		return nil, errors.Join(ErrGetScheme, err)
	}

	return NewApp(cfg, scheme, inputPath), nil
}

func getHybridScheme(inputPath string, msgSizeBytes uint64) (layer2.IScheme, error) {
	keyPath := filepath.Join(inputPath, hlk_settings.CPathKey)
	privKey, err := keys.GetPrivKey(keyPath)
	if err != nil {
		return nil, err
	}
	pubPath := filepath.Join(inputPath, hlk_settings.CPathPubKey)
	if _, err := keys.GetPubKey(privKey, pubPath); err != nil {
		return nil, err
	}
	return hybrid.NewScheme(privKey, msgSizeBytes)
}

func getSymmetricScheme(msgSizeBytes uint64) (layer2.IScheme, error) {
	return symmetric.NewScheme(msgSizeBytes)
}
