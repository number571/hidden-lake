package network

import (
	"context"
	"time"

	"github.com/number571/go-peer/pkg/anonymity"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	gopeer_logger "github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/pkg/adapters"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"
)

type IHiddenLakeNode interface {
	types.IRunner
	GetOriginNode() anonymity.INode

	SendRequest(context.Context, asymmetric.IPubKey, request.IRequest) error
	FetchRequest(context.Context, asymmetric.IPubKey, request.IRequest) (response.IResponse, error)
}

type ISettings interface {
	ISrvSettings
	IQBPSettings

	GetAdapterSettings() adapters.ISettings
}

type IQBPSettings interface {
	GetPowParallel() uint64
	GetQBPConsumers() uint64
	GetFetchTimeout() time.Duration
	GetQueuePeriod() time.Duration
	GetQueuePoolCap() [2]uint64
}

type ISrvSettings interface {
	GetLogger() gopeer_logger.ILogger
	GetServiceName() string
}
