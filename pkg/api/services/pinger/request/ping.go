package request

import (
	"net/http"

	hls_settings "github.com/number571/hidden-lake/internal/services/pinger/pkg/settings"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
)

func NewPingRequest() hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost(hls_settings.CAppShortName).
		WithPath(hls_settings.CPingPath).
		Build()
}
