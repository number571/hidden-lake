package adapters

import (
	"github.com/number571/go-peer/pkg/anonymity/qb/adapters"
	"github.com/number571/go-peer/pkg/types"
)

type IRunnerAdapter interface {
	types.IRunner
	adapters.IAdapter
}

type ISettings interface {
	GetNetworkKey() string
	GetWorkSizeBits() uint64
	GetMessageSizeBytes() uint64
}
