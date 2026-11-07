package limiter

import (
	"context"
	"net/http"

	hls_settings "github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/pldsize"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
	hlk_response "github.com/number571/hidden-lake/pkg/network/response"
)

var (
	gNotifyReqSize = uint64(len(
		hlk_request.NewRequestBuilder().
			WithMethod(http.MethodPost).
			WithPath(hls_settings.CAppShortName + "/notify").
			WithBody([]byte{}).
			Build().
			ToBytes(),
	))
	gNotifyRspSize = uint64(len(
		hlk_response.NewResponseBuilder().
			WithCode(200).
			WithBody([]byte{}).
			Build().
			ToBytes(),
	))
)

func GetLimitOnNotifyRequestSize(pCtx context.Context, pHlkClient hlk_client.IClient) (uint64, error) {
	return pldsize.GetLimitBySize(pCtx, pHlkClient, gNotifyReqSize)
}

func GetLimitOnNotifyResponseSize(pCtx context.Context, pHlkClient hlk_client.IClient) (uint64, error) {
	return pldsize.GetLimitBySize(pCtx, pHlkClient, gNotifyRspSize)
}
