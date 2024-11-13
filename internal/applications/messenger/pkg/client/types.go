package client

import (
	"context"

	hls_request "github.com/number571/hidden-lake/pkg/request"
)

type IClient interface {
	PingMessage(context.Context, string) error
	PushMessage(context.Context, string, []byte) error
}

type IRequester interface {
	PingMessage(context.Context, string, hls_request.IRequest) error
	PushMessage(context.Context, string, hls_request.IRequest) error
}

type IBuilder interface {
	PingMessage() hls_request.IRequest
	PushMessage([]byte) hls_request.IRequest
}
