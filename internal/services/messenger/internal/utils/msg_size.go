package utils

import (
	"context"
	"errors"

	"github.com/number571/hidden-lake/internal/kernel/pkg/client"
	hls_messenger_client "github.com/number571/hidden-lake/internal/services/messenger/pkg/client"
)

var (
	gReqSize = uint64(len(
		hls_messenger_client.NewBuilder().PushMessage([]byte{}).ToBytes(),
	))
)

func GetMessageLimit(pCtx context.Context, pHlsClient client.IClient) (uint64, error) {
	sett, err := pHlsClient.GetSettings(pCtx)
	if err != nil {
		return 0, errors.Join(ErrGetSettingsHLS, err)
	}

	pldLimit := sett.GetPayloadSizeBytes()
	if gReqSize >= pldLimit {
		return 0, ErrMessageSizeGteLimit
	}

	return pldLimit - gReqSize, nil
}
