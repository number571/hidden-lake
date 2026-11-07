package request

import (
	"fmt"
	"net/http"
	"net/url"

	hls_settings "github.com/number571/hidden-lake/internal/services/publisher/pkg/settings"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
)

func NewMsgboxPushRequest(pUID string, pBody []byte) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodPost).
		WithHost(hls_settings.CAppShortName).
		WithPath(fmt.Sprintf(
			"%s?uid=%s",
			hls_settings.CMsgboxPath,
			url.QueryEscape(pUID),
		)).
		WithBody(pBody).
		Build()
}

func NewMsgboxLoadRequestByUID(pUID string) hlk_request.IRequest {
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

func NewMsgboxLoadRequestByEnum(pEnum uint64) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost(hls_settings.CAppShortName).
		WithPath(fmt.Sprintf(
			"%s?nid=%d",
			hls_settings.CMsgboxPath,
			pEnum,
		)).
		WithBody([]byte{}).
		Build()
}
