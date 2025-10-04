package client

import (
	"context"

	hlk_request "github.com/number571/hidden-lake/pkg/request"
)

type IClient interface {
	Ping(context.Context, string) error
}

type IRequester interface {
	Ping(context.Context, string, hlk_request.IRequest) error
}

type IBuilder interface {
	Ping() hlk_request.IRequest
}
