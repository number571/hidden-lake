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
	GetAnonymityNode() anonymity.INode

	SendRequest(context.Context, asymmetric.IPubKey, request.IRequest) error
	FetchRequest(context.Context, asymmetric.IPubKey, request.IRequest) (response.IResponse, error)
}

type ISettings interface {
	ISubSettings
	GetAdapterSettings() adapters.ISettings
	GetQueuePeriod() time.Duration
	GetFetchTimeout() time.Duration
}

type ISubSettings interface {
	GetLogger() gopeer_logger.ILogger
	GetParallel() uint64
	GetServiceName() string
}
