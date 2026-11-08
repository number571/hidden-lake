package request

import (
	"fmt"
	"net/http"

	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
)

func NewListRequest(pPage uint64, pPublic bool) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost(hls_settings.CAppShortName).
		WithPath(fmt.Sprintf(
			"%s?page=%d&pub=%t",
			hls_settings.CListPath,
			pPage,
			pPublic,
		)).
		Build()
}
