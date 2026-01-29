package client

import (
	"context"
	"io"

	fileinfo "github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
)

type IClient interface {
	GetIndex(context.Context) (string, error)

	GetFileInfo(context.Context, string, string) (fileinfo.IFileInfo, error)
	GetListFiles(context.Context, string, uint64) (fileinfo.IFileInfoList, error)
	DownloadFile(io.Writer, context.Context, string, string) (bool, string, error)
}

type IRequester interface {
	GetIndex(context.Context) (string, error)

	GetFileInfo(context.Context, string, string) (fileinfo.IFileInfo, error)
	GetListFiles(context.Context, string, uint64) (fileinfo.IFileInfoList, error)
	DownloadFile(io.Writer, context.Context, string, string) (bool, string, error)
}
