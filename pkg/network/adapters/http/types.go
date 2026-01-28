package http

import (
	"net/http"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/pkg/network/adapters"
)

type IHTTPAdapter interface {
	adapters.IRunnerAdapter

	WithLogger(string, logger.ILogger) IHTTPAdapter
	WithHandlers(map[string]http.HandlerFunc) IHTTPAdapter
	GetOnlines() []string
}

type ISettings interface {
	GetAdapterSettings() adapters.ISettings
	GetAddress() string
	GetReadTimeout() time.Duration
	GetHandleTimeout() time.Duration
}
