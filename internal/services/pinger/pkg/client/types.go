package client

import (
	"context"

	hls_request "github.com/number571/hidden-lake/pkg/request"
)

type IClient interface {
	Ping(context.Context, string) error
}

type IRequester interface {
	Ping(context.Context, string, hls_request.IRequest) error
}

type IBuilder interface {
	Ping() hls_request.IRequest
}
