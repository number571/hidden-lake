package app

import (
	"errors"
	"path/filepath"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/anonymity/queue"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/internal/service/internal/handler"

	"github.com/number571/go-peer/pkg/client"
	net_message "github.com/number571/go-peer/pkg/network/message"
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

	p.fNode = anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FServiceName:  hls_settings.CServiceName,
			FNetworkMask:  hls_settings.CNetworkMask,
			FFetchTimeout: time.Duration(cfgSettings.GetFetchTimeoutMS()) * time.Millisecond,
		}),
		// Insecure to use logging in real anonymity projects!
		// Logging should only be used in overview or testing;
		p.fAnonLogger,
		kvDatabase,
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FAddress:      cfg.GetAddress().GetTCP(),
				FMaxConnects:  hls_settings.CNetworkMaxConns,
				FReadTimeout:  hls_settings.CNetworkReadTimeout,
				FWriteTimeout: hls_settings.CNetworkWriteTimeout,
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FMessageSettings:       cfgSettings,
					FLimitMessageSizeBytes: cfgSettings.GetMessageSizeBytes(),
					FWaitReadTimeout:       hls_settings.CConnWaitReadTimeout,
					FDialTimeout:           hls_settings.CConnDialTimeout,
					FReadTimeout:           hls_settings.CNetworkReadTimeout,
					FWriteTimeout:          hls_settings.CNetworkWriteTimeout,
				}),
			}),
			cache.NewLRUCache(hls_settings.CNetworkQueueCapacity),
		),
		queue.NewQBProblemProcessor(
			queue.NewSettings(&queue.SSettings{
				FMessageConstructSettings: net_message.NewConstructSettings(&net_message.SConstructSettings{
					FSettings: cfgSettings,
					FParallel: p.fParallel,
				}),
				FNetworkMask: hls_settings.CNetworkMask,
				FQueuePeriod: time.Duration(cfgSettings.GetQueuePeriodMS()) * time.Millisecond,
				FPoolCapacity: [2]uint64{
					hls_settings.CQueueMainPoolCapacity,
					hls_settings.CQueueRandPoolCapacity,
				},
			}),
			client,
		),
		func() asymmetric.IMapPubKeys {
			f2f := asymmetric.NewMapPubKeys()
			for _, pubKey := range cfg.GetFriends() {
				f2f.SetPubKey(pubKey)
			}
			return f2f
		}(),
	).HandleFunc(
		hls_settings.CServiceMask,
		handler.HandleServiceTCP(cfg),
	)

	return nil
}
