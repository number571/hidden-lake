package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/pkg/types"
	consumer_app "github.com/number571/hidden-lake/internal/adapters/chatingar/internal/consumer/pkg/app"
	producer_app "github.com/number571/hidden-lake/internal/adapters/chatingar/internal/producer/pkg/app"
	"github.com/number571/hidden-lake/internal/adapters/chatingar/pkg/app/config"
	"github.com/number571/hidden-lake/internal/adapters/chatingar/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
)

func InitApp(pArgs []string) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(flag.GetFlagValue(pArgs, "path", "."), "/")

	cfgPath := filepath.Join(inputPath, settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil, flag.GetFlagValue(pArgs, "network", ""))
	if err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}

	return NewApp(
		cfg,
		consumer_app.NewApp(cfg),
		producer_app.NewApp(cfg),
	), nil
}
