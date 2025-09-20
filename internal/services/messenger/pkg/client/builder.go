package client

import (
	"net/http"

	hls_messenger_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
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
		WithHost(hls_messenger_settings.CAppShortName).
		WithPath(hls_messenger_settings.CPushPath).
		WithBody(pBody).
		Build()
}
