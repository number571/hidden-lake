package utils

import (
	"context"
	"errors"

	hln_client "github.com/number571/hidden-lake/internal/applications/notifier/pkg/client"
	"github.com/number571/hidden-lake/internal/service/pkg/client"
)

var (
	gReqSize = uint64(len(
		hln_client.NewBuilder().Redirect(0, make([]byte, hln_client.CSaltSize), []byte{}).ToBytes(),
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
