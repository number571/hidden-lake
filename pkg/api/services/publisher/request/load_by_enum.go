package request

import (
	"fmt"
	"net/http"

	hls_settings "github.com/number571/hidden-lake/internal/services/publisher/pkg/settings"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
)

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
