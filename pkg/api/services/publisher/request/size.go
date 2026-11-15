package request

import (
	"fmt"
	"net/http"

	hls_settings "github.com/number571/hidden-lake/internal/services/publisher/pkg/settings"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
)

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
