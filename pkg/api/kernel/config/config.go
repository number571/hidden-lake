package config

import (
	"time"

	"github.com/number571/go-peer/pkg/crypto/hybrid/client"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/hidden-lake/internal/kernel/pkg/app/config"
)

func GetConfigSettings(pCfg config.IConfig, pClient client.IClient) SConfigSettings {
	sett := pCfg.GetSettings()
	return SConfigSettings{
		SConfigSettings: config.SConfigSettings{
			FNetworkKey:       sett.GetNetworkKey(),
			FMessageSizeBytes: sett.GetMessageSizeBytes(),
			FWorkSizeBits:     sett.GetWorkSizeBits(),
			FQBPConsumers:     sett.GetQBPConsumers(),
			FPowParallel:      sett.GetPowParallel(),
			FFetchTimeoutMS:   uint64(sett.GetFetchTimeout() / time.Millisecond), //nolint:gosec
			FQueuePeriodMS:    uint64(sett.GetQueuePeriod() / time.Millisecond),  //nolint:gosec
		},
		// encoding.CSizeUint64 = payload64.Head()
		FPayloadSizeBytes: pClient.GetPayloadLimit() - encoding.CSizeUint64,
	}
}
