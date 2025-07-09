package client

import (
	"context"

	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	hls_request "github.com/number571/hidden-lake/pkg/request"
)

type IClient interface {
	GetListFiles(context.Context, string, uint64) ([]hls_filesharer_settings.SFileInfo, error)
	LoadFileChunk(context.Context, string, string, uint64) ([]byte, error)
}

type IRequester interface {
	GetListFiles(context.Context, string, hls_request.IRequest) ([]hls_filesharer_settings.SFileInfo, error)
	LoadFileChunk(context.Context, string, hls_request.IRequest) ([]byte, error)
}

type IBuilder interface {
	GetListFiles(uint64) hls_request.IRequest
	LoadFileChunk(string, uint64) hls_request.IRequest
}
