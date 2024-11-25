package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/internal/helpers/loader/pkg/app/config"
	"github.com/number571/hidden-lake/internal/helpers/loader/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
)

func InitApp(pArgs []string) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(flag.GetStringFlagValue(pArgs, []string{"p", "path"}, "."), "/")

	cfgPath := filepath.Join(inputPath, settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil, flag.GetStringFlagValue(pArgs, []string{"n", "network"}, ""))
	if err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}

	return NewApp(cfg), nil
}
