package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/number571/hidden-lake/internal/helpers/traffic/internal/config"
	"github.com/number571/hidden-lake/internal/helpers/traffic/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"

	"github.com/number571/go-peer/pkg/types"
)

func InitApp(pArgs []string, pDefaultPath string) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(flag.GetFlagValue(pArgs, "path", pDefaultPath), "/")

	cfgPath := filepath.Join(inputPath, settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil)
	if err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}

	return NewApp(cfg, inputPath), nil
}
