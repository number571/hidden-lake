package limiter

import (
	"context"
	"net/http"

	hls_settings "github.com/number571/hidden-lake/internal/services/publisher/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/pldsize"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
	hlk_response "github.com/number571/hidden-lake/pkg/network/response"
)

var (
	gLoadReqSize = uint64(len(
		hlk_request.NewRequestBuilder().
			WithMethod(http.MethodGet).
			WithPath(hls_settings.CAppShortName + "/load?uid=0000000000000000000000&pub=false").
			// OR => WithPath(hls_settings.CAppShortName + "/load?nid=18446744073709551616&pub=false").
			WithBody([]byte{}).
			Build().
			ToBytes(),
	))
	gLoadRspSize = uint64(len(
		hlk_response.NewResponseBuilder().
			WithCode(200).
			WithHead(map[string]string{
				hls_settings.CHeaderMessageStateOnLoad:     "0",
				hls_settings.CHeaderAnotherMessageIDOnLoad: "0000000000000000000000",
			}).
			WithBody([]byte{}).
			Build().
			ToBytes(),
	))
)

func GetLimitOnLoadRequestSize(pCtx context.Context, pHlkClient hlk_client.IClient) (uint64, error) {
	return pldsize.GetLimitBySize(pCtx, pHlkClient, gLoadReqSize)
}

func GetLimitOnLoadResponseSize(pCtx context.Context, pHlkClient hlk_client.IClient) (uint64, error) {
	return pldsize.GetLimitBySize(pCtx, pHlkClient, gLoadRspSize)
}
