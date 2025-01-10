package client

import (
	"net/http"

	hlm_settings "github.com/number571/hidden-lake/internal/applications/messenger/pkg/settings"
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

func (p *sBuilder) PushMessage(pBody []byte) hls_request.IRequest {
	return hls_request.NewRequestBuilder().
		WithMethod(http.MethodPost).
		WithHost(hlm_settings.CServiceFullName).
		WithPath(hlm_settings.CPushPath).
		WithBody(pBody).
		Build()
}
