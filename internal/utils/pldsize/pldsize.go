package pldsize

import (
	"context"
	"errors"

	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
)

func GetLimitBySize(pCtx context.Context, pHlkClient hlk_client.IClient, pSize uint64) (uint64, error) {
	sett, err := pHlkClient.GetSettings(pCtx)
	if err != nil {
		return 0, errors.Join(ErrGetSettingsHLS, err)
	}

	pldLimit := sett.GetPayloadSizeBytes()
	if pSize >= pldLimit {
		return 0, ErrMessageSizeGteLimit
	}

	return pldLimit - pSize, nil
}
