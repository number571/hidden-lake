package app

import (
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/connkeeper"
	hiddenlake "github.com/number571/hidden-lake"
)

func (p *sApp) initConnKeeper(pNode network.INode) {
	p.fConnKeeper = connkeeper.NewConnKeeper(
		connkeeper.NewSettings(&connkeeper.SSettings{
			FConnections: func() []string { return p.fConfig.GetConnections() },
			FDuration:    hiddenlake.GSettings.GetKeeperPeriod(),
		}),
		pNode,
	)
}
