package client

import (
	"net/http"

	"github.com/number571/go-peer/pkg/message/layer1"
	hln_settings "github.com/number571/hidden-lake/internal/applications/notifier/pkg/settings"
	hls_request "github.com/number571/hidden-lake/pkg/request"
)

func init() {
	if len(hln_settings.CFinalyzePath) != len(hln_settings.CRedirectPath) {
		panic("notifier builder length")
	}
}

var (
	_ IBuilder = &sBuilder{}
)

type sBuilder struct {
}

func NewBuilder() IBuilder {
	return &sBuilder{}
}

func (p *sBuilder) Finalyze(pMsg layer1.IMessage) hls_request.IRequest {
	return basicBuilder(pMsg).
		WithPath(hln_settings.CFinalyzePath).
		Build()
}

func (p *sBuilder) Redirect(pMsg layer1.IMessage) hls_request.IRequest {
	return basicBuilder(pMsg).
		WithPath(hln_settings.CRedirectPath).
		Build()
}

func basicBuilder(pMsg layer1.IMessage) hls_request.IRequestBuilder {
	return hls_request.NewRequestBuilder().
		WithMethod(http.MethodPost).
		WithHost(hln_settings.CServiceFullName).
		WithBody(pMsg.ToBytes())
}
