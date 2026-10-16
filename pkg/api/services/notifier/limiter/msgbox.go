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
	gMsgboxReqSize = uint64(len(
		hlk_request.NewRequestBuilder().
			WithMethod(http.MethodPost).
			WithPath(hls_settings.CAppShortName + "/msgbox?uid=0000000000000000000000").
			WithBody([]byte{}).
			Build().
			ToBytes(),
	))
	gMsgboxRspSize = uint64(len(
		hlk_response.NewResponseBuilder().
			WithCode(200).
			WithBody([]byte{}).
			Build().
			ToBytes(),
	))
)

func GetLimitOnMsgboxRequestSize(pCtx context.Context, pHlkClient hlk_client.IClient) (uint64, error) {
	return pldsize.GetLimitBySize(pCtx, pHlkClient, gMsgboxReqSize)
}

func GetLimitOnMsgboxResponseSize(pCtx context.Context, pHlkClient hlk_client.IClient) (uint64, error) {
	return pldsize.GetLimitBySize(pCtx, pHlkClient, gMsgboxRspSize)
}
