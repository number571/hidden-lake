package request

import (
	"net/http"

	hls_settings "github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
)

func NewNotifyRequest(pBody []byte) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodPost).
		WithHost(hls_settings.CAppShortName).
		WithPath(hls_settings.CNotifyPath).
		WithBody(pBody).
		Build()
}
