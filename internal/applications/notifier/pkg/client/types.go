package client

import (
	"context"

	"github.com/number571/go-peer/pkg/message/layer1"
	hls_request "github.com/number571/hidden-lake/pkg/request"
)

const (
	CSaltSize = 32
)

type IClient interface {
	Finalyze(context.Context, []string, layer1.IMessage) error
	Redirect(context.Context, []string, string, layer1.IMessage) error
}

type IRequester interface {
	Broadcast(context.Context, []string, hls_request.IRequest) error
}

type IBuilder interface {
	Finalyze(layer1.IMessage) hls_request.IRequest
	Redirect(layer1.IMessage) hls_request.IRequest
}

type ISettings interface {
	GetWorkSizeBits() uint64
	GetPowParallel() uint64
}
