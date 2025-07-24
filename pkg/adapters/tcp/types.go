package tcp

import (
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/connkeeper"
	"github.com/number571/hidden-lake/pkg/adapters"
	"github.com/number571/hidden-lake/pkg/utils/appname"
)

type ITCPAdapter interface {
	adapters.IRunnerAdapter

	WithLogger(appname.IFmtAppName, logger.ILogger) ITCPAdapter
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
