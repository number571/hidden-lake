package network

import (
	"context"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	gopeer_logger "github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"
	gopeer_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"
)

type IHandlerF func(
	context.Context,
	asymmetric.IPubKey,
	request.IRequest,
) (response.IResponse, error)

type IHiddenLakeNode interface {
	types.IRunner
	GetOriginNode() anonymity.INode

	SendRequest(context.Context, asymmetric.IPubKey, request.IRequest) error
	FetchRequest(context.Context, asymmetric.IPubKey, request.IRequest) (response.IResponse, error)
}

type ISettings interface {
	ISubSettings
	GetMessageSettings() gopeer_message.ISettings
	GetMessageSizeBytes() uint64
	GetQueuePeriod() time.Duration
	GetFetchTimeout() time.Duration
}

type ISubSettings interface {
	GetLogger() gopeer_logger.ILogger
	GetParallel() uint64
	GetTCPAddress() string
	GetServiceName() string
}
