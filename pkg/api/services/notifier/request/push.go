package request

import (
	"fmt"
	"net/http"
	"net/url"

	hls_settings "github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
)

func NewNotifyRequest(pBody string) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodPost).
		WithHost(hls_settings.CAppShortName).
		WithPath(hls_settings.CNotifyPath).
		WithBody([]byte(pBody)).
		Build()
}

func NewPushRequest(pUID, pBody string) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodPost).
		WithHost(hls_settings.CAppShortName).
		WithPath(fmt.Sprintf(
			"%s?uid=%s",
			hls_settings.CPushPath,
			url.QueryEscape(pUID),
		)).
		WithBody([]byte(pBody)).
		Build()
}
