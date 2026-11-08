package request

import (
	"fmt"
	"net/http"
	"net/url"

	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
)

func NewInfoRequest(pFileName string, pPublic bool) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost(hls_settings.CAppShortName).
		WithPath(fmt.Sprintf(
			"%s?name=%s&pub=%t",
			hls_settings.CInfoPath,
			url.QueryEscape(pFileName),
			pPublic,
		)).
		Build()
}
