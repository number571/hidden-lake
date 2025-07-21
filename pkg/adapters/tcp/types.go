package tcp

import (
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/connkeeper"
	"github.com/number571/hidden-lake/internal/utils/name"
	"github.com/number571/hidden-lake/pkg/adapters"
)

type ITCPAdapter interface {
	adapters.IRunnerAdapter

	WithLogger(name.IAppName, logger.ILogger) ITCPAdapter
	GetConnKeeper() connkeeper.IConnKeeper
}

type ISettings interface {
	ISrvSettings

	GetAdapterSettings() adapters.ISettings
}

type ISrvSettings interface {
	GetAddress() string
	GetConnNumLimit() uint64
	GetConnKeepPeriod() time.Duration
	GetSendTimeout() time.Duration
	GetRecvTimeout() time.Duration
	GetDialTimeout() time.Duration
	GetWaitTimeout() time.Duration
}
