package utils

import (
	"context"
	"errors"

	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	hlk_response "github.com/number571/hidden-lake/pkg/response"
)

var (
	gRespSize = uint64(len(
		hlk_response.NewResponseBuilder().
			WithCode(200).
			WithHead(map[string]string{
				"Content-Type":                   api.CApplicationOctetStream,
				hlk_settings.CHeaderResponseMode: hlk_settings.CHeaderResponseModeON,
			}).
			WithBody([]byte{}).
			Build().
			ToBytes(),
	))
)

func GetMessageLimitOnLoadPage(pCtx context.Context, pHlkClient hlk_client.IClient) (uint64, error) {
	sett, err := pHlkClient.GetSettings(pCtx)
	if err != nil {
		return 0, errors.Join(ErrGetSettingsHLS, err)
	}

	pldLimit := sett.GetPayloadSizeBytes()
	if gRespSize >= pldLimit {
		return 0, ErrMessageSizeGteLimit
	}

	return pldLimit - gRespSize, nil
}
