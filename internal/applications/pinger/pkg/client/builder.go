package client

import (
	"net/http"

	hlp_settings "github.com/number571/hidden-lake/internal/applications/pinger/pkg/settings"
	hls_request "github.com/number571/hidden-lake/pkg/request"
)

var (
	_ IBuilder = &sBuilder{}
)

type sBuilder struct {
}

func NewBuilder() IBuilder {
	return &sBuilder{}
}

func (p *sBuilder) Ping() hls_request.IRequest {
	return hls_request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost(hlp_settings.CServiceFullName).
		WithPath(hlp_settings.CPingPath).
		Build()
}
