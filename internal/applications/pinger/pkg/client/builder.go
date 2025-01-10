package client

import (
	"net/http"

	hlr_settings "github.com/number571/hidden-lake/internal/applications/pinger/pkg/settings"
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
		WithHost(hlr_settings.CServiceFullName).
		WithPath(hlr_settings.CPingPath).
		Build()
}
