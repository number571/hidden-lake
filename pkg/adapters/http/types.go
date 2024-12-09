package http

import (
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/utils/name"
	"github.com/number571/hidden-lake/pkg/adapters"
)

type IHTTPAdapter interface {
	WithLogger(name.IServiceName, logger.ILogger) IHTTPAdapter
	WithHandlers(...IHandler) IHTTPAdapter
	GetOnlines() []string
	adapters.IRunnerAdapter
}

type IHandler interface {
	GetPath() string
	GetFunc() func(http.ResponseWriter, *http.Request)
}

type ISettings interface {
	GetAdapterSettings() adapters.ISettings
	GetAddress() string
}
