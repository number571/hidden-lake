package tcp

import (
	"github.com/number571/go-peer/pkg/network/connkeeper"
	"github.com/number571/hidden-lake/pkg/adapters"
)

type ITCPAdapter interface {
	GetConnKeeper() connkeeper.IConnKeeper
	adapters.IRunnerAdapter
}

type ISettings interface {
	GetAdapterSettings() adapters.ISettings
	GetAddress() string
}
