package client

import (
	"context"
	"io"

	"github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
)

type IClient interface {
	GetIndex(context.Context) error

	GetRemoteList(context.Context, string, uint64, bool) ([]dto.IFileInfo, error)
	GetRemoteFile(io.Writer, context.Context, string, string, bool) (bool, error)
	DelRemoteFile(context.Context, string, string, bool) error
	GetRemoteFileInfo(context.Context, string, string, bool) (dto.IFileInfo, error)

	GetRemoteFileProc(context.Context, string, string, bool) (dto.IDownloadProcess, error)
	DelRemoteFileProc(context.Context, string, string, bool) error
	GetRemoteListProc(context.Context) ([]dto.IDownloadProcess, error)

	GetLocalList(context.Context, string, uint64) ([]dto.IFileInfo, error)
	GetLocalFile(io.Writer, context.Context, string, string) error
	PutLocalFile(context.Context, string, string, io.Reader) error
	DelLocalFile(context.Context, string, string) error
	GetLocalFileInfo(context.Context, string, string) (dto.IFileInfo, error)
}

type IRequester interface {
	GetIndex(context.Context) error

	GetRemoteList(context.Context, string, uint64, bool) ([]dto.IFileInfo, error)
	GetRemoteFile(io.Writer, context.Context, string, string, bool) (bool, error)
	DelRemoteFile(context.Context, string, string, bool) error
	GetRemoteFileInfo(context.Context, string, string, bool) (dto.IFileInfo, error)

	GetRemoteFileProc(context.Context, string, string, bool) (dto.IDownloadProcess, error)
	DelRemoteFileProc(context.Context, string, string, bool) error
	GetRemoteListProc(context.Context) ([]dto.IDownloadProcess, error)

	GetLocalList(context.Context, string, uint64) ([]dto.IFileInfo, error)
	GetLocalFile(io.Writer, context.Context, string, string) error
	PutLocalFile(context.Context, string, string, io.Reader) error
	DelLocalFile(context.Context, string, string) error
	GetLocalFileInfo(context.Context, string, string) (dto.IFileInfo, error)
}
