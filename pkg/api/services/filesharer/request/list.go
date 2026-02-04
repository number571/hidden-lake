package request

import (
	"fmt"
	"net/http"

	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
)

func NewListRequest(pPage uint64, pPersonal bool) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost(hls_settings.CAppShortName).
		WithPath(fmt.Sprintf(
			"%s?page=%d&personal=%t",
			hls_settings.CListPath,
			pPage,
			pPersonal,
		)).
		Build()
}
