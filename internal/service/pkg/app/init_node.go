package app

import (
	"errors"
	"path/filepath"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/pkg/adapters"
	"github.com/number571/hidden-lake/pkg/adapters/http"
	"github.com/number571/hidden-lake/pkg/network"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/hidden-lake/internal/service/internal/handler"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

func (p *sApp) initAnonNode() error {
	var (
		cfg         = p.fCfgW.GetConfig()
		cfgSettings = cfg.GetSettings()
	)

	client := client.NewClient(p.fPrivKey, cfgSettings.GetMessageSizeBytes())
	if client.GetPayloadLimit() <= encoding.CSizeUint64 {
		return ErrMessageSizeLimit
	}

	kvDatabase, err := database.NewKVDatabase(filepath.Join(p.fPathTo, hls_settings.CPathDB))
	if err != nil {
		return errors.Join(ErrOpenKVDatabase, err)
	}

	adapterSettings := adapters.NewSettings(&adapters.SSettings{
		FWorkSizeBits:     cfgSettings.GetWorkSizeBits(),
		FNetworkKey:       cfgSettings.GetNetworkKey(),
		FMessageSizeBytes: cfgSettings.GetMessageSizeBytes(),
	})

	node := network.NewHiddenLakeNode(
		network.NewSettings(&network.SSettings{
			FAdapterSettings: adapterSettings,
			FQBPSettings: &network.SQBPSettings{
				FQueuePeriod:  cfgSettings.GetQueuePeriod(),
				FFetchTimeout: cfgSettings.GetFetchTimeout(),
				FPowParallel:  cfgSettings.GetPowParallel(),
				FQBPConsumers: cfgSettings.GetQBPConsumers(),
				FQueuePoolCap: cfgSettings.GetQueuePoolCap(),
			},
			FSrvSettings: &network.SSrvSettings{
				FServiceName: hls_settings.GServiceName.Short(),
				FLogger:      p.fAnonLogger,
			},
		}),
		p.fPrivKey,
		kvDatabase,
		http.NewHTTPAdapter(
			http.NewSettings(&http.SSettings{
				FAdapterSettings: adapterSettings,
				FAddress:         cfg.GetAddress().GetExternal(),
			}),
			cache.NewLRUCache(build.GetSettings().FNetworkManager.FCacheHashesCap),
			func() []string { return p.fCfgW.GetConfig().GetEndpoints() },
		),
		handler.HandleServiceFunc(cfg, p.fAnonLogger),
	)

	originNode := node.GetOriginNode()
	for _, f := range cfg.GetFriends() {
		originNode.GetMapPubKeys().SetPubKey(f)
	}

	p.fNode = node
	return nil
}
