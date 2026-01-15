package client

import (
	"context"
	"io"

	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/utils"
)

type IClient interface {
	GetIndex(context.Context) (string, error)

	GetFileInfo(context.Context, string, string) (utils.IFileInfo, error)
	GetListFiles(context.Context, string, uint64) ([]utils.IFileInfo, error)
	DownloadFile(io.Writer, context.Context, string, string) error
}

type IRequester interface {
	GetIndex(context.Context) (string, error)

	GetFileInfo(context.Context, string, string) (utils.IFileInfo, error)
	GetListFiles(context.Context, string, uint64) ([]utils.IFileInfo, error)
	DownloadFile(io.Writer, context.Context, string, string) error
}
