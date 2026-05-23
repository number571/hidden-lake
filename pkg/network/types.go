package network

import (
	"context"
	"time"

	anonymity "github.com/number571/go-peer/pkg/anonymity/qb"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	gopeer_logger "github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/pkg/network/adapters"
	"github.com/number571/hidden-lake/pkg/network/request"
	"github.com/number571/hidden-lake/pkg/network/response"
)

type IHiddenLakeNode interface {
	types.IRunner
	GetOriginNode() anonymity.INode

	SendRequest(context.Context, layer2.IParticipantKey, request.IRequest) error
	FetchRequest(context.Context, layer2.IParticipantKey, request.IRequest) (response.IResponse, error)
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
}

type ISrvSettings interface {
	GetLogger() gopeer_logger.ILogger
	GetFmtAppName() string
}
