package request

import (
	"fmt"
	"net/http"
	"net/url"

	hls_settings "github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
)

func NewMsgboxPushRequest(pUID, pBody string) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodPost).
		WithHost(hls_settings.CAppShortName).
		WithPath(fmt.Sprintf(
			"%s?uid=%s",
			hls_settings.CMsgboxPath,
			url.QueryEscape(pUID),
		)).
		WithBody([]byte(pBody)).
		Build()
}

func NewMsgboxLoadRequest(pUID string) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost(hls_settings.CAppShortName).
		WithPath(fmt.Sprintf(
			"%s?uid=%s",
			hls_settings.CMsgboxPath,
			url.QueryEscape(pUID),
		)).
		WithBody([]byte{}).
		Build()
}
