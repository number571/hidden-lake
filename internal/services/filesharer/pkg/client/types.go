package client

import (
	"context"
	"io"

	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/client/fileinfo"
)

type IClient interface {
	GetIndex(context.Context) (string, error)

	GetFileInfo(context.Context, string, string) (fileinfo.IFileInfo, error)
	GetListFiles(context.Context, string, uint64) ([]fileinfo.IFileInfo, error)
	DownloadFile(io.Writer, context.Context, string, string) error
}

type IRequester interface {
	GetIndex(context.Context) (string, error)

	GetFileInfo(context.Context, string, string) (fileinfo.IFileInfo, error)
	GetListFiles(context.Context, string, uint64) ([]fileinfo.IFileInfo, error)
	DownloadFile(io.Writer, context.Context, string, string) error
}
