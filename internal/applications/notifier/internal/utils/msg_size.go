package utils

import (
	"context"
	"errors"

	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/payload"
	hln_client "github.com/number571/hidden-lake/internal/applications/notifier/pkg/client"
	"github.com/number571/hidden-lake/internal/service/pkg/client"
)

var (
	gReqSize = uint64(len(
		hln_client.NewBuilder().Redirect(layer1.NewMessage(
			layer1.NewConstructSettings(&layer1.SConstructSettings{
				FSettings: layer1.NewSettings(&layer1.SSettings{}),
			}),
			payload.NewPayload32(0, []byte{}),
		)).ToBytes(),
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
