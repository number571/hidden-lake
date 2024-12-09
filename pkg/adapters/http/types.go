package http

import (
	"net/http"

	"github.com/number571/hidden-lake/pkg/adapters"
)

type IHTTPAdapter interface {
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
