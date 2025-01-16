package config

import (
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/hidden-lake/internal/service/pkg/app/config"
)

func GetConfigSettings(pCfg config.IConfig, pClient client.IClient) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FNetworkKey:       sett.GetNetworkKey(),
			FMessageSizeBytes: sett.GetMessageSizeBytes(),
			FWorkSizeBits:     sett.GetWorkSizeBits(),
			FQBConsumers:      sett.GetQBConsumers(),
			FPowParallel:      sett.GetPowParallel(),
			FFetchTimeoutMS:   uint64(sett.GetFetchTimeout() / time.Millisecond),
			FQueuePeriodMS:    uint64(sett.GetQueuePeriod() / time.Millisecond),
		},
		// encoding.CSizeUint64 = payload64.Head()
		FPayloadSizeBytes: pClient.GetPayloadLimit() - encoding.CSizeUint64,
	}
}
