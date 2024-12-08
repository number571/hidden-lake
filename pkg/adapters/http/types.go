package http

import "github.com/number571/hidden-lake/pkg/adapters"

type IHTTPAdapter interface {
	adapters.IRunnerAdapter
}

type ISettings interface {
	GetAddress() string
	GetNetworkKey() string
	GetWorkSizeBits() uint64
	GetMessageSizeBytes() uint64
}
