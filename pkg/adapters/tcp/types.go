package tcp

import (
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/connkeeper"
	"github.com/number571/hidden-lake/internal/utils/name"
	"github.com/number571/hidden-lake/pkg/adapters"
)

type ITCPAdapter interface {
	WithLogger(name.IServiceName, logger.ILogger) ITCPAdapter
	GetConnKeeper() connkeeper.IConnKeeper
	adapters.IRunnerAdapter
}

type ISettings interface {
	GetAdapterSettings() adapters.ISettings
	GetAddress() string
}
