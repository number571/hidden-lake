package http

import "github.com/number571/hidden-lake/pkg/adapters"

const (
	CHandleNetworkAdapterPath = "/api/network/adapter"
)

type IHTTPAdapter interface {
	GetOnlines() []string
	adapters.IRunnerAdapter
}

type ISettings interface {
	GetAddress() string
	GetProducePath() string
	GetNetworkKey() string
	GetWorkSizeBits() uint64
	GetMessageSizeBytes() uint64
}
