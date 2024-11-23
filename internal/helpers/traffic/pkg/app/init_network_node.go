package app

import (
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/helpers/traffic/internal/cache"
	"github.com/number571/hidden-lake/internal/helpers/traffic/internal/handler"
	"github.com/number571/hidden-lake/internal/helpers/traffic/internal/storage"
)

func (p *sApp) initNetworkNode(pStorage storage.IMessageStorage) {
	cfgSettings := p.fConfig.GetSettings()
	p.fNode = network.NewNode(
		network.NewSettings(&network.SSettings{
			FAddress:      p.fConfig.GetAddress().GetTCP(),
			FMaxConnects:  build.GSettings.FNetworkManager.FConnectsLimiter,
			FReadTimeout:  build.GSettings.GetReadTimeout(),
			FWriteTimeout: build.GSettings.GetWriteTimeout(),
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FMessageSettings:       cfgSettings,
				FLimitMessageSizeBytes: cfgSettings.GetMessageSizeBytes(),
				FWaitReadTimeout:       build.GSettings.GetWaitTimeout(),
				FDialTimeout:           build.GSettings.GetDialTimeout(),
				FReadTimeout:           build.GSettings.GetReadTimeout(),
				FWriteTimeout:          build.GSettings.GetWriteTimeout(),
			}),
		}),
		cache.NewLRUCache(build.GSettings.FNetworkManager.FCacheHashesCap),
	).HandleFunc(
		build.GSettings.FProtoMask.FNetwork,
		handler.HandleServiceTCP(p.fConfig, pStorage, p.fAnonLogger),
	)
}
