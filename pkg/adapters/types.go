package adapters

import (
	"github.com/number571/go-peer/pkg/anonymity/adapters"
	"github.com/number571/go-peer/pkg/types"
)

type IRunnerAdapter interface {
	types.IRunner
	adapters.IAdapter
}
