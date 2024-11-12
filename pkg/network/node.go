package network

import (
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/anonymity/queue"
	"github.com/number571/go-peer/pkg/network/conn"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	hiddenlake "github.com/number571/hidden-lake"
)

func NewHiddenLakeNode(
	pSettings ISettings,
	pPrivKey asymmetric.IPrivKey,
	pKVDatabase database.IKVDatabase,
) anonymity.INode {
	return anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FNetworkMask:  hiddenlake.GSettings.FProtoMask.FNetwork,
			FServiceName:  pSettings.GetServiceName(),
			FFetchTimeout: pSettings.GetFetchTimeout(),
		}),
		// Insecure to use logging in real anonymity projects!
		// Logging should only be used in overview or testing;
		pSettings.GetLogger(),
		pKVDatabase,
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FAddress:      pSettings.GetTCPAddress(),
				FMaxConnects:  hiddenlake.GSettings.FNetworkManager.FConnectsLimiter,
				FReadTimeout:  hiddenlake.GSettings.GetReadTimeout(),
				FWriteTimeout: hiddenlake.GSettings.GetWriteTimeout(),
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FMessageSettings:       pSettings.GetMessageSettings(),
					FLimitMessageSizeBytes: pSettings.GetMessageSizeBytes(),
					FWaitReadTimeout:       hiddenlake.GSettings.GetWaitTimeout(),
					FDialTimeout:           hiddenlake.GSettings.GetDialTimeout(),
					FReadTimeout:           hiddenlake.GSettings.GetReadTimeout(),
					FWriteTimeout:          hiddenlake.GSettings.GetWriteTimeout(),
				}),
			}),
			cache.NewLRUCache(hiddenlake.GSettings.FNetworkManager.FCacheHashesCap),
		),
		queue.NewQBProblemProcessor(
			queue.NewSettings(&queue.SSettings{
				FMessageConstructSettings: net_message.NewConstructSettings(&net_message.SConstructSettings{
					FSettings: pSettings.GetMessageSettings(),
					FParallel: pSettings.GetParallel(),
				}),
				FNetworkMask: hiddenlake.GSettings.FProtoMask.FNetwork,
				FQueuePeriod: pSettings.GetQueuePeriod(),
				FPoolCapacity: [2]uint64{
					hiddenlake.GSettings.FQueueCapacity.FMain,
					hiddenlake.GSettings.FQueueCapacity.FRand,
				},
			}),
			func() client.IClient {
				client := client.NewClient(pPrivKey, pSettings.GetMessageSizeBytes())
				if client.GetPayloadLimit() <= encoding.CSizeUint64 {
					panic(`client.GetPayloadLimit() <= encoding.CSizeUint64`)
				}
				return client
			}(),
		),
		asymmetric.NewMapPubKeys(),
	)
}
