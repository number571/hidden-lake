package client

import (
	"context"
	"io"

	"github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
)

type IClient interface {
	GetIndex(pCtx context.Context) error

	GetRemoteList(pCtx context.Context, pFriend string, pPage uint64) (uint64, []string, error)
	GetRemoteFile(pW io.Writer, pCtx context.Context, pFriend string, pFileName string) (bool, error)
	DelRemoteFile(pCtx context.Context, pFriend string, pFileName string) error
	GetRemoteFileInfo(pCtx context.Context, pFriend string, pFileName string) (dto.IFileInfo, error)

	GetRemoteFileProc(pCtx context.Context, pFriend string, pFileName string) (dto.IDownloadProcess, error)
	DelRemoteFileProc(pCtx context.Context, pFriend string, pFileName string) error
	GetRemoteListProc(pCtx context.Context) ([]dto.IDownloadProcess, error)

	GetLocalList(pCtx context.Context, pFriend string, pPage uint64) (uint64, []string, error)
	GetLocalFile(pW io.Writer, pCtx context.Context, pFriend string, pFileName string) error
	PutLocalFile(pCtx context.Context, pFriend string, pFileName string, pR io.Reader) error
	DelLocalFile(pCtx context.Context, pFriend string, pFileName string) error
	GetLocalFileInfo(pCtx context.Context, pFriend string, pFileName string) (dto.IFileInfo, error)
}
