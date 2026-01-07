package client

import (
	"context"

	hlk_request "github.com/number571/hidden-lake/pkg/request"
)

type IFileInfo interface {
	GetName() string
	GetHash() string
	GetSize() uint64
}

type IClient interface {
	GetFileInfo(context.Context, string, string) (IFileInfo, error)
	GetListFiles(context.Context, string, uint64) ([]IFileInfo, error)
	LoadFileChunk(context.Context, string, string, uint64) ([]byte, error)
}

type IRequester interface {
	GetFileInfo(context.Context, string, hlk_request.IRequest) (IFileInfo, error)
	GetListFiles(context.Context, string, hlk_request.IRequest) ([]IFileInfo, error)
	LoadFileChunk(context.Context, string, hlk_request.IRequest) ([]byte, error)
}

type IBuilder interface {
	GetFileInfo(pFileName string) hlk_request.IRequest
	GetListFiles(uint64) hlk_request.IRequest
	LoadFileChunk(string, uint64) hlk_request.IRequest
}
