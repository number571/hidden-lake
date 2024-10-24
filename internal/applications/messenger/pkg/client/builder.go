package client

import (
	"net/http"

	hlm_settings "github.com/number571/hidden-lake/internal/applications/messenger/pkg/settings"
	hls_request "github.com/number571/hidden-lake/internal/service/pkg/request"
)

var (
	_ IBuilder = &sBuilder{}
)

type sBuilder struct {
}

func NewBuilder() IBuilder {
	return &sBuilder{}
}

func (p *sBuilder) PingMessage() hls_request.IRequest {
	return hls_request.NewRequest(
		http.MethodGet,
		hlm_settings.CServiceFullName,
		hlm_settings.CPingPath,
	)
}

func (p *sBuilder) PushMessage(pBody []byte) hls_request.IRequest {
	return hls_request.NewRequest(
		http.MethodPost,
		hlm_settings.CServiceFullName,
		hlm_settings.CPushPath,
	).WithBody(pBody)
}
