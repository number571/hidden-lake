package http

import (
	"net/http"

	"github.com/number571/hidden-lake/pkg/adapters"
)

type IHTTPAdapter interface {
	WithHandlers(...IHandlerFunc) IHTTPAdapter
	GetOnlines() []string
	adapters.IRunnerAdapter
}

type IHandlerFunc interface {
	GetPath() string
	GetFunc() func(http.ResponseWriter, *http.Request)
}

type ISettings interface {
	GetAdapterSettings() adapters.ISettings
	GetAddress() string
}
