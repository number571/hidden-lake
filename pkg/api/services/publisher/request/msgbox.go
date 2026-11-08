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
			hls_settings.CPushPath,
			url.QueryEscape(pUID),
		)).
		WithBody(pBody).
		Build()
}

func NewMsgboxLoadRequestByUID(pUID string, pPublic bool) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost(hls_settings.CAppShortName).
		WithPath(fmt.Sprintf(
			"%s?uid=%s&pub=%t",
			hls_settings.CLoadPath,
			url.QueryEscape(pUID),
			pPublic,
		)).
		WithBody([]byte{}).
		Build()
}

func NewMsgboxLoadRequestByEnum(pEnum uint64, pPublic bool) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost(hls_settings.CAppShortName).
		WithPath(fmt.Sprintf(
			"%s?nid=%d&pub=%t",
			hls_settings.CLoadPath,
			pEnum,
			pPublic,
		)).
		WithBody([]byte{}).
		Build()
}

func NewMsgboxSizeRequest(pPublic bool) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost(hls_settings.CAppShortName).
		WithPath(hls_settings.CSizePath).
		WithPath(fmt.Sprintf(
			"%s?pub=%t",
			hls_settings.CSizePath,
			pPublic,
		)).
		WithBody([]byte{}).
		Build()
}
