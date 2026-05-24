package app

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/pkg/network"
	"github.com/number571/hidden-lake/pkg/network/adapters"
	"github.com/number571/hidden-lake/pkg/network/adapters/http"

	"github.com/number571/hidden-lake/internal/kernel/internal/handler"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
)

func (p *sApp) initAnonNode() error {
	var (
		cfg           = p.fCfgW.GetConfig()
		cfgSettings   = cfg.GetSettings()
		buildSettings = build.GetSettings()
	)

	kvDatabase, err := database.NewKVDatabase(filepath.Join(p.fPathTo, hlk_settings.CPathDB))
	if err != nil {
		return errors.Join(ErrOpenKVDatabase, err)
	}

	adapterSettings := adapters.NewSettings(&adapters.SSettings{
		FWorkSizeBits:     cfgSettings.GetWorkSizeBits(),
		FNetworkKey:       cfgSettings.GetNetworkKey(),
		FMessageSizeBytes: cfgSettings.GetMessageSizeBytes(),
	})

	node, err := network.NewHiddenLakeNode(
		network.NewSettings(&network.SSettings{
			FAdapterSettings: adapterSettings,
			FQBPSettings: &network.SQBPSettings{
				FQueuePeriod:  cfgSettings.GetQueuePeriod(),
				FFetchTimeout: cfgSettings.GetFetchTimeout(),
				FPowParallel:  cfgSettings.GetPowParallel(),
				FQBPConsumers: cfgSettings.GetQBPConsumers(),
			},
			FServeSettings: &network.SServeSettings{
				FServiceName: hlk_settings.GetAppShortNameFMT(),
				FLogger:      p.fAnonLogger,
			},
		}),
		p.fScheme,
		layer2.NewKeysContainer(),
		kvDatabase,
		http.NewHTTPAdapter(
			http.NewSettings(&http.SSettings{
				FAdapterSettings: adapterSettings,
				FServeSettings: &http.SServeSettings{
					FAddress:       cfg.GetAddress().GetExternal(),
					FSubscribeID:   fmt.Sprintf("%s-%s", hlk_settings.CAppShortName, random.NewRandom().GetString(16)),
					FReadTimeout:   buildSettings.GetHttpReadTimeout(),
					FHandleTimeout: buildSettings.GetHttpHandleTimeout(),
					FDataBrokerParams: [2]uint64{
						buildSettings.FNetworkManager.FMessageBuffer,
						buildSettings.FNetworkManager.FConnNumLimit,
					},
				},
			}),
			cache.NewLRUCache(buildSettings.FStorageManager.FCacheHashesCap),
			func() []string { return p.fCfgW.GetConfig().GetEndpoints() },
		),
		handler.HandleServiceFunc(cfg, p.fAnonLogger),
	)
	if err != nil {
		return errors.Join(ErrCreateNode, err)
	}

	for _, k := range cfg.GetFriends() {
		if ok := node.GetOriginNode().GetKeysContainer().Add(k); !ok {
			return ErrAddFriendToList
		}
	}

	p.fNode = node
	return nil
}
