package client

import (
	"context"

	hls_request "github.com/number571/hidden-lake/pkg/request"
)

const (
	CSaltSize = 32
)

type IClient interface {
	Initialize(context.Context, []string, []byte) ([]byte, error)
	Finalyze(context.Context, []string, uint64, []byte, []byte) error
	Redirect(context.Context, []string, string, uint64, []byte, []byte) error
}

type IRequester interface {
	Broadcast(context.Context, []string, hls_request.IRequest) error
}

type IBuilder interface {
	Finalyze(uint64, []byte, []byte) hls_request.IRequest
	Redirect(uint64, []byte, []byte) hls_request.IRequest
}

type ISettings interface {
	GetDiffBits() uint64
	GetParallel() uint64
}
