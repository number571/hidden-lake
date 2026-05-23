package config

import (
	"time"

	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/hidden-lake/internal/kernel/pkg/app/config"
)

func GetConfigSettings(pCfg config.IConfig, pScheme layer2.IScheme) SConfigSettings {
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
		FPayloadSizeBytes: pScheme.GetPayloadLimit() - encoding.CSizeUint64,
	}
}
