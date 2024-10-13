package client

import (
	"net/http"
	"strings"

	hlr_settings "github.com/number571/hidden-lake/cmd/applications/remoter/pkg/settings"
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

func (p *sBuilder) Exec(pCmd ...string) hls_request.IRequest {
	return hls_request.NewRequest(
		http.MethodPost,
		hlr_settings.CServiceFullName,
		hlr_settings.CExecPath,
	).WithBody([]byte(strings.Join(pCmd, " ")))
}
