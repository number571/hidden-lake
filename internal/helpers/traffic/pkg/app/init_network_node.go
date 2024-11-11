package app

import (
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	hiddenlake "github.com/number571/hidden-lake"
	"github.com/number571/hidden-lake/internal/helpers/traffic/internal/cache"
	"github.com/number571/hidden-lake/internal/helpers/traffic/internal/handler"
	"github.com/number571/hidden-lake/internal/helpers/traffic/internal/storage"
)

func (p *sApp) initNetworkNode(pStorage storage.IMessageStorage) {
	cfgSettings := p.fConfig.GetSettings()
	p.fNode = network.NewNode(
		network.NewSettings(&network.SSettings{
			FAddress:      p.fConfig.GetAddress().GetTCP(),
			FMaxConnects:  hiddenlake.GSettings.FNetworkManager.FConnectsLimiter,
			FReadTimeout:  hiddenlake.GSettings.GetReadTimeout(),
			FWriteTimeout: hiddenlake.GSettings.GetWriteTimeout(),
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FMessageSettings:       cfgSettings,
				FLimitMessageSizeBytes: cfgSettings.GetMessageSizeBytes(),
				FWaitReadTimeout:       hiddenlake.GSettings.GetWaitTimeout(),
				FDialTimeout:           hiddenlake.GSettings.GetDialTimeout(),
				FReadTimeout:           hiddenlake.GSettings.GetReadTimeout(),
				FWriteTimeout:          hiddenlake.GSettings.GetWriteTimeout(),
			}),
		}),
		cache.NewLRUCache(hiddenlake.GSettings.FNetworkManager.FCacheHashesCap),
	).HandleFunc(
		hiddenlake.GSettings.FProtoMask.FService,
		handler.HandleServiceTCP(p.fConfig, pStorage, p.fAnonLogger),
	)
}
