package client

import (
	"context"

	hls_settings "github.com/number571/hidden-lake/internal/services/remoter/pkg/settings"
)

type IClient interface {
	GetIndex(context.Context) (string, error)
	ExecCommand(context.Context, string, ...string) ([]byte, error)
}

type IRequester interface {
	GetIndex(context.Context) (string, error)
	ExecCommand(context.Context, string, *hls_settings.SCommandExecRequest) ([]byte, error)
}

type IBuilder interface {
	ExecCommand(...string) *hls_settings.SCommandExecRequest
}
