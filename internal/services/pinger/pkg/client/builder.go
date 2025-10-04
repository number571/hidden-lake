package client

import (
	"net/http"

	hls_pinger_settings "github.com/number571/hidden-lake/internal/services/pinger/pkg/settings"
	hlk_request "github.com/number571/hidden-lake/pkg/request"
)

var (
	_ IBuilder = &sBuilder{}
)

type sBuilder struct {
}

func NewBuilder() IBuilder {
	return &sBuilder{}
}

func (p *sBuilder) Ping() hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost(hls_pinger_settings.CAppShortName).
		WithPath(hls_pinger_settings.CPingPath).
		Build()
}
