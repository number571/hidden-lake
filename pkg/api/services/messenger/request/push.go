package request

import (
	"net/http"

	hls_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
)

func NewPushRequest(pBody string) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodPost).
		WithHost(hls_settings.CAppShortName).
		WithPath(hls_settings.CPushPath).
		WithBody([]byte(pBody)).
		Build()
}
