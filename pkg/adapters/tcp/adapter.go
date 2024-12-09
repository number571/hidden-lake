package tcp

import (
	"context"
	"errors"
	"sync"

	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/connkeeper"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/hidden-lake/build"

	net_message "github.com/number571/go-peer/pkg/network/message"
)

const (
	netMessageChanSize = 32
)

var (
	_ ITCPAdapter = &sTCPAdapter{}
)

type sTCPAdapter struct {
	fNetMsgChan chan net_message.IMessage
	fConnKeeper connkeeper.IConnKeeper
}

func NewTCPAdapter(
	pSettings ISettings,
	pCache cache.ICache,
	pConnsGetter func() []string,
) ITCPAdapter {
	adapterSettings := pSettings.GetAdapterSettings()
	tcpAdapter := &sTCPAdapter{
		fNetMsgChan: make(chan net_message.IMessage, netMessageChanSize),
		fConnKeeper: connkeeper.NewConnKeeper(
			connkeeper.NewSettings(&connkeeper.SSettings{
				FDuration:    build.GSettings.GetKeeperPeriod(),
				FConnections: pConnsGetter,
			}),
			network.NewNode(
				network.NewSettings(&network.SSettings{
					FAddress:      pSettings.GetAddress(),
					FMaxConnects:  build.GSettings.FNetworkManager.FConnectsLimiter,
					FReadTimeout:  build.GSettings.GetReadTimeout(),
					FWriteTimeout: build.GSettings.GetWriteTimeout(),
					FConnSettings: conn.NewSettings(&conn.SSettings{
						FMessageSettings:       adapterSettings,
						FLimitMessageSizeBytes: adapterSettings.GetMessageSizeBytes(),
						FWaitReadTimeout:       build.GSettings.GetWaitTimeout(),
						FDialTimeout:           build.GSettings.GetDialTimeout(),
						FReadTimeout:           build.GSettings.GetReadTimeout(),
						FWriteTimeout:          build.GSettings.GetWriteTimeout(),
					}),
				}),
				pCache,
			),
		),
	}
	tcpAdapter.fConnKeeper.GetNetworkNode().HandleFunc(
		build.GSettings.FProtoMask.FNetwork,
		func(_ context.Context, _ network.INode, _ conn.IConn, msg net_message.IMessage) error {
			tcpAdapter.fNetMsgChan <- msg
			return nil
		},
	)
	return tcpAdapter
}

func (p *sTCPAdapter) GetConnKeeper() connkeeper.IConnKeeper {
	return p.fConnKeeper
}

func (p *sTCPAdapter) Run(pCtx context.Context) error {
	chCtx, cancel := context.WithCancel(pCtx)
	defer cancel()

	const N = 2

	errs := make([]error, N)
	wg := &sync.WaitGroup{}
	wg.Add(N)

	go func() {
		defer func() { wg.Done(); cancel() }()
		errs[0] = p.fConnKeeper.GetNetworkNode().Run(chCtx)
	}()

	go func() {
		defer func() { wg.Done(); cancel() }()
		errs[1] = p.fConnKeeper.Run(chCtx)
	}()

	wg.Wait()

	select {
	case <-pCtx.Done():
		return pCtx.Err()
	default:
		return errors.Join(errs...)
	}
}

func (p *sTCPAdapter) Produce(pCtx context.Context, pNetMsg net_message.IMessage) error {
	networkNode := p.fConnKeeper.GetNetworkNode()
	return networkNode.BroadcastMessage(pCtx, pNetMsg)
}

func (p *sTCPAdapter) Consume(pCtx context.Context) (net_message.IMessage, error) {
	select {
	case <-pCtx.Done():
		return nil, pCtx.Err()
	case msg := <-p.fNetMsgChan:
		return msg, nil
	}
}
