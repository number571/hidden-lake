package client

import (
	"net/http"

	hlm_settings "github.com/number571/hidden-lake/cmd/applications/messenger/pkg/settings"
	hls_request "github.com/number571/hidden-lake/cmd/service/pkg/request"
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
	return hls_request.NewRequest(
		http.MethodPost,
		hlm_settings.CServiceFullName,
		hlm_settings.CPushPath,
	).WithBody(pBody)
}
