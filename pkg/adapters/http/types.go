package http

import (
	"net/http"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/pkg/adapters"
	"github.com/number571/hidden-lake/pkg/utils/appname"
)

type IHTTPAdapter interface {
	adapters.IRunnerAdapter

	WithLogger(appname.IAppName, logger.ILogger) IHTTPAdapter
	WithHandlers(map[string]http.HandlerFunc) IHTTPAdapter
	GetOnlines() []string
}

type ISettings interface {
	GetAdapterSettings() adapters.ISettings
	GetAddress() string
	GetReadTimeout() time.Duration
	GetHandleTimeout() time.Duration
}
