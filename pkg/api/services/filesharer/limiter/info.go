package limiter

import (
	"context"
	"net/http"

	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/pldsize"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
	hlk_response "github.com/number571/hidden-lake/pkg/network/response"
)

var (
	gInfoReqSize = uint64(len(
		hlk_request.NewRequestBuilder().
			WithMethod(http.MethodGet).
			WithPath(hls_settings.CAppShortName + "/info?name=&personal=false").
			WithBody([]byte{}).
			Build().
			ToBytes(),
	))
	gInfoRspSize = uint64(len(
		hlk_response.NewResponseBuilder().
			WithCode(200).
			WithBody([]byte(`{"name":"","size":0,"hash":"0000000000000000000000000000000000000000"}`)).
			Build().
			ToBytes(),
	))
)

func GetLimitOnInfoRequestSize(pCtx context.Context, pHlkClient hlk_client.IClient) (uint64, error) {
	return pldsize.GetLimitBySize(pCtx, pHlkClient, gInfoReqSize)
}

func GetLimitOnInfoResponseSize(pCtx context.Context, pHlkClient hlk_client.IClient) (uint64, error) {
	return pldsize.GetLimitBySize(pCtx, pHlkClient, gInfoRspSize)
}
