package app

import (
	"errors"
	"path/filepath"
	"time"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/storage/database"
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

	kvDatabase, err := database.NewKVDatabase(filepath.Join(p.fPathTo, hls_settings.CPathDB))
	if err != nil {
		return errors.Join(ErrOpenKVDatabase, err)
	}

	client := client.NewClient(p.fPrivKey, cfgSettings.GetMessageSizeBytes())
	if client.GetPayloadLimit() <= encoding.CSizeUint64 {
		return errors.Join(ErrMessageSizeLimit, err)
	}

	p.fNode = network.NewHiddenLakeNode(
		network.NewSettings(&network.SSettings{
			FMessageSettings:  cfg.GetSettings(),
			FMessageSizeBytes: cfgSettings.GetMessageSizeBytes(),
			FQueuePeriod:      time.Duration(cfgSettings.GetQueuePeriodMS()) * time.Millisecond,
			FFetchTimeout:     time.Duration(cfgSettings.GetFetchTimeoutMS()) * time.Millisecond,
			SSubSettings: network.SSubSettings{
				FServiceName: hls_settings.CServiceName,
				FTCPAddress:  cfg.GetAddress().GetTCP(),
				FParallel:    p.fParallel,
				FLogger:      p.fAnonLogger,
			},
		}),
		p.fPrivKey,
		kvDatabase,
		func() []string { return p.fCfgW.GetConfig().GetConnections() },
		handler.HandleServiceTCP(cfg, p.fAnonLogger),
	)

	originNode := p.fNode.GetOrigNode()
	for _, f := range cfg.GetFriends() {
		originNode.GetMapPubKeys().SetPubKey(f)
	}

	return nil
}
