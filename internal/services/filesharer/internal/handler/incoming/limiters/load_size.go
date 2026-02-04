package limiters

import (
	"context"
	"errors"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	hlk_response "github.com/number571/hidden-lake/pkg/network/response"
)

var (
	gLoadRspSize = uint64(len(
		hlk_response.NewResponseBuilder().
			WithCode(200).
			WithHead(map[string]string{
				"Content-Type":               api.CApplicationOctetStream,
				hls_settings.CHeaderFileHash: hashing.NewHasher([]byte{}).ToString(),
			}).
			WithBody([]byte{}).
			Build().
			ToBytes(),
	))
)

func GetLimitOnLoadResponseSize(pCtx context.Context, pHlkClient hlk_client.IClient) (uint64, error) {
	sett, err := pHlkClient.GetSettings(pCtx)
	if err != nil {
		return 0, errors.Join(ErrGetSettingsHLS, err)
	}

	pldLimit := sett.GetPayloadSizeBytes()
	if gLoadRspSize >= pldLimit {
		return 0, ErrMessageSizeGteLimit
	}

	return pldLimit - gLoadRspSize, nil
}
