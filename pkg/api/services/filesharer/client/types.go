package client

import (
	"context"
	"io"

	fileinfo "github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
)

type IClient interface {
	GetIndex(context.Context) (string, error)

	GetRemoteList(context.Context, string, uint64, bool) (fileinfo.IFileInfoList, error)
	GetRemoteFile(io.Writer, context.Context, string, string, bool) (bool, error)
	DelRemoteFile(context.Context, string, string, bool) error
	GetRemoteFileInfo(context.Context, string, string, bool) (fileinfo.IFileInfo, error)

	GetLocalList(context.Context, string, uint64) (fileinfo.IFileInfoList, error)
	GetLocalFile(io.Writer, context.Context, string, string) error
	PutLocalFile(context.Context, string, string, io.Reader) error
	DelLocalFile(context.Context, string, string) error
	GetLocalFileInfo(context.Context, string, string) (fileinfo.IFileInfo, error)
}

type IRequester interface {
	GetIndex(context.Context) (string, error)

	GetRemoteFileInfo(context.Context, string, string, bool) (fileinfo.IFileInfo, error)
	GetRemoteFile(io.Writer, context.Context, string, string, bool) (bool, error)
	DelRemoteFile(context.Context, string, string, bool) error
	GetRemoteList(context.Context, string, uint64, bool) (fileinfo.IFileInfoList, error)

	GetLocalFileInfo(context.Context, string, string) (fileinfo.IFileInfo, error)
	GetLocalFile(io.Writer, context.Context, string, string) error
	PutLocalFile(context.Context, string, string, io.Reader) error
	DelLocalFile(context.Context, string, string) error
	GetLocalList(context.Context, string, uint64) (fileinfo.IFileInfoList, error)
}
