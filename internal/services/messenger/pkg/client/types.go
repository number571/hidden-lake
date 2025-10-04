package client

import (
	"context"

	hlk_request "github.com/number571/hidden-lake/pkg/request"
)

type IClient interface {
	PushMessage(context.Context, string, []byte) error
}

type IRequester interface {
	PushMessage(context.Context, string, hlk_request.IRequest) error
}

type IBuilder interface {
	PushMessage([]byte) hlk_request.IRequest
}
