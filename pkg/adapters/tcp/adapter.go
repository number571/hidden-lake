package tcp

import (
	"context"
	"errors"
	"sync"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/network/connkeeper"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/utils/name"

	anon_logger "github.com/number571/go-peer/pkg/anonymity/logger"
	"github.com/number571/go-peer/pkg/message/layer1"
	internal_anon_logger "github.com/number571/hidden-lake/internal/utils/logger/anon"
)

const (
	netMessageChanSize = 32
)

var (
	_ ITCPAdapter = &sTCPAdapter{}
)

type sTCPAdapter struct {
	fNetMsgChan chan layer1.IMessage
	fConnKeeper connkeeper.IConnKeeper

	fShortName string
	fLogger    logger.ILogger
}

func NewTCPAdapter(
	pSettings ISettings,
	pCache cache.ICache,
	pConnsGetter func() []string,
) ITCPAdapter {
	adapterSettings := pSettings.GetAdapterSettings()
	p := &sTCPAdapter{
		fNetMsgChan: make(chan layer1.IMessage, netMessageChanSize),
		fConnKeeper: connkeeper.NewConnKeeper(
			connkeeper.NewSettings(&connkeeper.SSettings{
				FDuration:    pSettings.GetConnKeepPeriod(),
				FConnections: pConnsGetter,
			}),
			network.NewNode(
				network.NewSettings(&network.SSettings{
					FAddress:      pSettings.GetAddress(),
					FMaxConnects:  pSettings.GetConnNumLimit(),
					FReadTimeout:  pSettings.GetRecvTimeout(),
					FWriteTimeout: pSettings.GetSendTimeout(),
					FConnSettings: conn.NewSettings(&conn.SSettings{
						FMessageSettings:       adapterSettings,
						FLimitMessageSizeBytes: adapterSettings.GetMessageSizeBytes(),
						FWaitReadTimeout:       pSettings.GetWaitTimeout(),
						FDialTimeout:           pSettings.GetDialTimeout(),
						FReadTimeout:           pSettings.GetRecvTimeout(),
						FWriteTimeout:          pSettings.GetSendTimeout(),
					}),
				}),
				pCache,
			),
		),
		fLogger: logger.NewLogger(
			logger.NewSettings(&logger.SSettings{}),
			func(_ logger.ILogArg) string { return "" },
		),
	}
	p.fConnKeeper.GetNetworkNode().HandleFunc(
		build.GetSettings().FProtoMask.FNetwork,
		func(_ context.Context, _ network.INode, conn conn.IConn, msg layer1.IMessage) error {
			logBuilder := anon_logger.NewLogBuilder(p.fShortName)
			p.fLogger.PushInfo(logBuilder.
				WithType(internal_anon_logger.CLogInfoRecvNetworkMessage).
				WithHash(msg.GetHash()).
				WithProof(msg.GetProof()).
				WithSize(len(msg.ToBytes())).
				WithConn(conn.GetSocket().RemoteAddr().String()))
			p.fNetMsgChan <- msg
			return nil
		},
	)
	return p
}

func (p *sTCPAdapter) WithLogger(pName name.IServiceName, pLogger logger.ILogger) ITCPAdapter {
	p.fShortName = pName.Short()
	p.fLogger = pLogger
	return p
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
		errs := append([]error{ErrRunning}, errs...)
		return errors.Join(errs...)
	}
}

func (p *sTCPAdapter) Produce(pCtx context.Context, pNetMsg layer1.IMessage) error {
	logBuilder := anon_logger.NewLogBuilder(p.fShortName)
	logBuilder.
		WithType(internal_anon_logger.CLogBaseSendNetworkMessage).
		WithHash(pNetMsg.GetHash()).
		WithProof(pNetMsg.GetProof()).
		WithSize(len(pNetMsg.ToBytes())).
		WithConn("tcp")

	networkNode := p.fConnKeeper.GetNetworkNode()
	if err := networkNode.BroadcastMessage(pCtx, pNetMsg); err != nil {
		if errors.Is(err, network.ErrNoConnections) {
			p.fLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnNoConnections))
		} else {
			p.fLogger.PushInfo(logBuilder)
		}
		return errors.Join(ErrBroadcast, err)
	}
	p.fLogger.PushInfo(logBuilder)
	return nil
}

func (p *sTCPAdapter) Consume(pCtx context.Context) (layer1.IMessage, error) {
	select {
	case <-pCtx.Done():
		return nil, pCtx.Err()
	case msg := <-p.fNetMsgChan:
		return msg, nil
	}
}
