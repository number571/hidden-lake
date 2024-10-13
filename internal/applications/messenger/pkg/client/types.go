package client

import (
	"context"

	hls_request "github.com/number571/hidden-lake/internal/service/pkg/request"
)

type IClient interface {
	PushMessage(context.Context, string, []byte) error
}

type IRequester interface {
	PushMessage(context.Context, string, hls_request.IRequest) error
}

type IBuilder interface {
	PushMessage([]byte) hls_request.IRequest
}
