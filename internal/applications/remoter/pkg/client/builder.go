package client

import (
	"net/http"
	"strings"

	hlr_settings "github.com/number571/hidden-lake/internal/applications/remoter/pkg/settings"
	hls_request "github.com/number571/hidden-lake/pkg/request"
)

var (
	_ IBuilder = &sBuilder{}
)

type sBuilder struct {
	fPassword string
}

func NewBuilder(pPassword string) IBuilder {
	return &sBuilder{
		fPassword: pPassword,
	}
}

func (p *sBuilder) Exec(pCmd ...string) hls_request.IRequest {
	return hls_request.NewRequestBuilder().
		WithMethod(http.MethodPost).
		WithHost(hlr_settings.CServiceFullName).
		WithPath(hlr_settings.CExecPath).
		WithHead(map[string]string{
			hlr_settings.CHeaderPassword: p.fPassword,
		}).
		WithBody([]byte(strings.Join(pCmd, hlr_settings.CExecSeparator))).
		Build()
}
