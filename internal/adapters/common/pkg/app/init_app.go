package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/pkg/types"
	consumer_app "github.com/number571/hidden-lake/internal/adapters/common/internal/consumer/pkg/app"
	producer_app "github.com/number571/hidden-lake/internal/adapters/common/internal/producer/pkg/app"
	"github.com/number571/hidden-lake/internal/adapters/common/pkg/app/config"
	"github.com/number571/hidden-lake/internal/adapters/common/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
)

func InitApp(pArgs []string, pFlags flag.IFlags) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(pFlags.Get("path").GetStringValue(pArgs), "/")

	cfgPath := filepath.Join(inputPath, settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil, pFlags.Get("network").GetStringValue(pArgs))
	if err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}

	return NewApp(
		cfg,
		consumer_app.NewApp(cfg),
		producer_app.NewApp(cfg),
	), nil
}
