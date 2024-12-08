package app

import (
	"errors"
	"path/filepath"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/pkg/adapters/http"
	"github.com/number571/hidden-lake/pkg/network"

	"github.com/number571/go-peer/pkg/client"
	hla_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
	"github.com/number571/hidden-lake/internal/service/internal/handler"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

func (p *sApp) initAnonNode() error {
	var (
		cfg         = p.fCfgW.GetConfig()
		cfgSettings = cfg.GetSettings()
	)

	kvDatabase, err := database.NewKVDatabase(filepath.Join(p.fPathTo, hls_settings.CPathDB))
	if err != nil {
		return errors.Join(ErrOpenKVDatabase, err)
	}

	client := client.NewClient(p.fPrivKey, cfgSettings.GetMessageSizeBytes())
	if client.GetPayloadLimit() <= encoding.CSizeUint64 {
		return errors.Join(ErrMessageSizeLimit, err)
	}

	settings := network.NewSettings(&network.SSettings{
		FWorkSizeBits:     cfgSettings.GetWorkSizeBits(),
		FNetworkKey:       cfgSettings.GetNetworkKey(),
		FMessageSizeBytes: cfgSettings.GetMessageSizeBytes(),
		FQueuePeriod:      cfgSettings.GetQueuePeriod(),
		FFetchTimeout:     cfgSettings.GetFetchTimeout(),
		FSubSettings: &network.SSubSettings{
			FServiceName: hls_settings.GServiceName.Short(),
			FParallel:    p.fParallel,
			FLogger:      p.fAnonLogger,
		},
	})

	node := network.NewHiddenLakeNode(
		settings,
		p.fPrivKey,
		kvDatabase,
		http.NewHTTPAdapter(
			http.NewSettings(&http.SSettings{
				FProducePath:      hla_settings.CHandleNetworkAdapterPath, // TODO: delete
				FAddress:          cfg.GetAddress().GetExternal(),
				FWorkSizeBits:     cfgSettings.GetWorkSizeBits(),
				FNetworkKey:       cfgSettings.GetNetworkKey(),
				FMessageSizeBytes: cfgSettings.GetMessageSizeBytes(),
			}),
			func() []string { return p.fCfgW.GetConfig().GetAdapters() },
		),
		handler.HandleServiceTCP(cfg, p.fAnonLogger),
	)

	originNode := node.GetAnonymityNode()
	for _, f := range cfg.GetFriends() {
		originNode.GetMapPubKeys().SetPubKey(f)
	}

	p.fNode = node
	return nil
}
