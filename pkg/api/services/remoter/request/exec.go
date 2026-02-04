package request

import (
	"net/http"
	"strings"

	hls_settings "github.com/number571/hidden-lake/internal/services/remoter/pkg/settings"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
)

func NewExecRequest(pPassword string, pCommand []string) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodPost).
		WithHost(hls_settings.CAppShortName).
		WithPath(hls_settings.CExecPath).
		WithHead(map[string]string{
			hls_settings.CHeaderPassword: pPassword,
		}).
		WithBody([]byte(strings.Join(pCommand, hls_settings.CExecSeparator))).
		Build()
}
