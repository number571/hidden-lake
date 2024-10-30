package utils

import (
	"context"
	"errors"

	hlm_client "github.com/number571/hidden-lake/internal/applications/messenger/pkg/client"
	"github.com/number571/hidden-lake/internal/service/pkg/client"
)

var (
	gReqSize = uint64(len(
		hlm_client.NewBuilder().PushMessage([]byte{}).ToBytes(),
	))
)

func GetMessageLimit(pCtx context.Context, pHlsClient client.IClient) (uint64, error) {
	sett, err := pHlsClient.GetSettings(pCtx)
	if err != nil {
		return 0, errors.Join(ErrGetSettingsHLS, err)
	}

	msgLimitOrig := sett.GetLimitMessageSizeBytes()
	if gReqSize >= msgLimitOrig {
		return 0, ErrMessageSizeGteLimit
	}

	return msgLimitOrig - gReqSize, nil
}
