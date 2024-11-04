package config

import (
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/hidden-lake/internal/service/pkg/app/config"
)

func GetConfigSettings(pCfg config.IConfig, pClient client.IClient) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FMessageSizeBytes: sett.GetMessageSizeBytes(),
			FWorkSizeBits:     sett.GetWorkSizeBits(),
			FFetchTimeoutMS:   sett.GetFetchTimeoutMS(),
			FQueuePeriodMS:    sett.GetQueuePeriodMS(),
			FNetworkKey:       sett.GetNetworkKey(),
		},
		// encoding.CSizeUint64 = payload64.Head()
		FLimitMessageSizeBytes: pClient.GetPayloadLimit() - encoding.CSizeUint64,
	}
}
