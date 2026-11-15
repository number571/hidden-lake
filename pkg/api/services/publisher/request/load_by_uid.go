package request

import (
	"fmt"
	"net/http"
	"net/url"

	hls_settings "github.com/number571/hidden-lake/internal/services/publisher/pkg/settings"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
)

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
