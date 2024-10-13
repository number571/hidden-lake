package utils

import (
	"context"

	"github.com/number571/go-peer/pkg/utils"
	hlm_client "github.com/number571/hidden-lake/cmd/applications/messenger/pkg/client"
	"github.com/number571/hidden-lake/cmd/service/pkg/client"
)

var (
	gReqSize = uint64(len(
		hlm_client.NewBuilder().PushMessage([]byte{}).ToBytes(),
	))
)

func GetMessageLimit(pCtx context.Context, pHlsClient client.IClient) (uint64, error) {
	sett, err := pHlsClient.GetSettings(pCtx)
	if err != nil {
		return 0, utils.MergeErrors(ErrGetSettingsHLS, err)
	}

	msgLimitOrig := sett.GetLimitMessageSizeBytes()
	if gReqSize >= msgLimitOrig {
		return 0, ErrMessageSizeGteLimit
	}

	return msgLimitOrig - gReqSize, nil
}
