package request

import (
	"fmt"
	"net/http"
	"net/url"

	hls_settings "github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
	"github.com/number571/hidden-lake/pkg/api/services/notifier/client/dto"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
)

func NewPushRequest(pMsg dto.IMessage) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodPost).
		WithHost(hls_settings.CAppShortName).
		WithPath(hls_settings.CPushPath).
		WithBody([]byte(pMsg.GetBody())).
		Build()
}

func NewPushRequestWithUID(pMsg dto.IMessage) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodPost).
		WithHost(hls_settings.CAppShortName).
		WithPath(fmt.Sprintf(
			"%s?uid=%s",
			hls_settings.CPushPath,
			url.QueryEscape(pMsg.GetUID()),
		)).
		WithBody([]byte(pMsg.GetBody())).
		Build()
}
