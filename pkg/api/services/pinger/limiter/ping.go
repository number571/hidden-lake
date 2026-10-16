package limiter

import (
	"context"
	"net/http"

	hls_settings "github.com/number571/hidden-lake/internal/services/pinger/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/pldsize"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
	hlk_response "github.com/number571/hidden-lake/pkg/network/response"
)

var (
	gPingReqSize = uint64(len(
		hlk_request.NewRequestBuilder().
			WithMethod(http.MethodGet).
			WithPath(hls_settings.CAppShortName + "/ping").
			WithBody([]byte{}).
			Build().
			ToBytes(),
	))
	gPingRspSize = uint64(len(
		hlk_response.NewResponseBuilder().
			WithCode(200).
			WithBody([]byte{}).
			Build().
			ToBytes(),
	))
)

func GetLimitOnPingRequestSize(pCtx context.Context, pHlkClient hlk_client.IClient) (uint64, error) {
	return pldsize.GetLimitBySize(pCtx, pHlkClient, gPingReqSize)
}

func GetLimitOnPingResponseSize(pCtx context.Context, pHlkClient hlk_client.IClient) (uint64, error) {
	return pldsize.GetLimitBySize(pCtx, pHlkClient, gPingRspSize)
}
