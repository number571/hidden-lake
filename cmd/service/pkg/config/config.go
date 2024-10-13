package config

import (
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/hidden-lake/cmd/service/internal/config"
)

func GetConfigSettings(pCfg config.IConfig, pClient client.IClient) SConfigSettings {
	sett := pCfg.GetSettings()
	msgLimit := pClient.GetMessageLimit()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FMessageSizeBytes:     sett.GetMessageSizeBytes(),
			FKeySizeBits:          sett.GetKeySizeBits(),
			FWorkSizeBits:         sett.GetWorkSizeBits(),
			FFetchTimeoutMS:       sett.GetFetchTimeoutMS(),
			FQueuePeriodMS:        sett.GetQueuePeriodMS(),
			FRandQueuePeriodMS:    sett.GetRandQueuePeriodMS(),
			FRandMessageSizeBytes: sett.GetRandMessageSizeBytes(),
			FNetworkKey:           sett.GetNetworkKey(),
		},
		// encoding.CSizeUint64 = payload64.Head()
		FLimitMessageSizeBytes: msgLimit - encoding.CSizeUint64,
	}
}
