package request

import (
	"fmt"
	"net/http"
	"net/url"

	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
)

func NewLoadRequest(pFileName string, pChunk uint64, pPersonal bool) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost(hls_settings.CAppShortName).
		WithPath(fmt.Sprintf(
			"%s?name=%s&chunk=%d&personal=%t",
			hls_settings.CLoadPath,
			url.QueryEscape(pFileName),
			pChunk,
			pPersonal,
		)).
		Build()
}
