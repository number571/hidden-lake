package client

import (
	"context"

	hlk_request "github.com/number571/hidden-lake/pkg/request"
)

type IClient interface {
	Exec(context.Context, string, ...string) ([]byte, error)
}

type IRequester interface {
	Exec(context.Context, string, hlk_request.IRequest) ([]byte, error)
}

type IBuilder interface {
	Exec(...string) hlk_request.IRequest
}
