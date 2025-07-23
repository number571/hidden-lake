package client

import (
	"net/http"
	"strings"

	hls_remoter_settings "github.com/number571/hidden-lake/internal/services/remoter/pkg/settings"
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
		WithHost(hls_remoter_settings.CAppFullName).
		WithPath(hls_remoter_settings.CExecPath).
		WithHead(map[string]string{
			hls_remoter_settings.CHeaderPassword: p.fPassword,
		}).
		WithBody([]byte(strings.Join(pCmd, hls_remoter_settings.CExecSeparator))).
		Build()
}
