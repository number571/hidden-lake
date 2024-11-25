package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/internal/helpers/traffic/pkg/app/config"
	"github.com/number571/hidden-lake/internal/helpers/traffic/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
)

func InitApp(pArgs []string, pFlags flag.IFlags) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(pFlags.Get("path").GetStringValue(pArgs), "/")

	cfgPath := filepath.Join(inputPath, settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil, pFlags.Get("network").GetStringValue(pArgs))
	if err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}

	return NewApp(cfg, inputPath), nil
}
