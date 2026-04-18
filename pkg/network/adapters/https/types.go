package https

import (
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/pkg/network/adapters"
)

type IHTTPSAdapter interface {
	adapters.IRunnerAdapter

	WithLogger(string, logger.ILogger) IHTTPSAdapter

	GetOnlines() []string
}

type ISettings interface {
	ISrvSettings

	GetAdapterSettings() adapters.ISettings
}

type ISrvSettings interface {
	GetAddress() string
	GetRateLimitParams() [2]float64
	GetDataBrokerParam() uint64
	GetAuthMapper() map[string]string
	GetReadTimeout() time.Duration
	GetHandleTimeout() time.Duration
}
