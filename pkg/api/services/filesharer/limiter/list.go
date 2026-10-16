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
	gListReqSize = uint64(len(
		hlk_request.NewRequestBuilder().
			WithMethod(http.MethodGet).
			WithPath(hls_settings.CAppShortName + "/list?page=&personal=false").
			WithBody([]byte{}).
			Build().
			ToBytes(),
	))
	gListRspSize = uint64(len(
		hlk_response.NewResponseBuilder().
			WithCode(200).
			WithBody([]byte{}).
			Build().
			ToBytes(),
	))
)

func GetLimitOnListRequestSize(pCtx context.Context, pHlkClient hlk_client.IClient) (uint64, error) {
	return pldsize.GetLimitBySize(pCtx, pHlkClient, gListReqSize)
}

func GetLimitOnListResponseSize(pCtx context.Context, pHlkClient hlk_client.IClient) (uint64, error) {
	return pldsize.GetLimitBySize(pCtx, pHlkClient, gListRspSize)
}
