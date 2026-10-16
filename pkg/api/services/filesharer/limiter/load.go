package limiter

import (
	"context"
	"net/http"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/pldsize"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
	hlk_response "github.com/number571/hidden-lake/pkg/network/response"
)

var (
	gLoadReqSize = uint64(len(
		hlk_request.NewRequestBuilder().
			WithMethod(http.MethodGet).
			WithPath(hls_settings.CAppShortName + "?name=&chunk=&personal=false").
			WithBody([]byte{}).
			Build().
			ToBytes(),
	))
	gLoadRspSize = uint64(len(
		hlk_response.NewResponseBuilder().
			WithCode(200).
			WithHead(map[string]string{
				hls_settings.CHeaderFileHashOnLoad: hashing.NewHasher([]byte{}).ToString()[:hls_settings.CTrimHashSizeOnLoad],
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
