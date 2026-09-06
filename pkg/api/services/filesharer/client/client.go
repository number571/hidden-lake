package client

import (
	"context"
	"io"

	"github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
)

var (
	_ IClient = &sClient{}
)

type sClient struct {
	fRequester IRequester
}

func NewClient(pRequester IRequester) IClient {
	return &sClient{
		fRequester: pRequester,
	}
}

func (p *sClient) GetIndex(pCtx context.Context) error {
	return p.fRequester.GetIndex(pCtx)
}

func (p *sClient) GetRemoteList(pCtx context.Context, pFriend string, pPage uint64, pPersonal bool) ([]dto.IFileInfo, error) {
	return p.fRequester.GetRemoteList(pCtx, pFriend, pPage, pPersonal)
}

func (p *sClient) GetRemoteFile(pW io.Writer, pCtx context.Context, pFriend string, pFilename string, pPersonal bool) (bool, error) {
	return p.fRequester.GetRemoteFile(pW, pCtx, pFriend, pFilename, pPersonal)
}

func (p *sClient) DelRemoteFile(pCtx context.Context, pFriend string, pFilename string, pPersonal bool) error {
	return p.fRequester.DelRemoteFile(pCtx, pFriend, pFilename, pPersonal)
}

func (p *sClient) GetRemoteFileInfo(pCtx context.Context, pFriend string, pFilename string, pPersonal bool) (dto.IFileInfo, error) {
	return p.fRequester.GetRemoteFileInfo(pCtx, pFriend, pFilename, pPersonal)
}

func (p *sClient) GetLocalList(pCtx context.Context, pFriend string, pPage uint64) ([]dto.IFileInfo, error) {
	return p.fRequester.GetLocalList(pCtx, pFriend, pPage)
}

func (p *sClient) GetLocalFile(pW io.Writer, pCtx context.Context, pFriend string, pFilename string) error {
	return p.fRequester.GetLocalFile(pW, pCtx, pFriend, pFilename)
}

func (p *sClient) PutLocalFile(pCtx context.Context, pFriend string, pFilename string, pR io.Reader) error {
	return p.fRequester.PutLocalFile(pCtx, pFriend, pFilename, pR)
}

func (p *sClient) DelLocalFile(pCtx context.Context, pFriend string, pFilename string) error {
	return p.fRequester.DelLocalFile(pCtx, pFriend, pFilename)
}

func (p *sClient) GetLocalFileInfo(pCtx context.Context, pFriend string, pFilename string) (dto.IFileInfo, error) {
	return p.fRequester.GetLocalFileInfo(pCtx, pFriend, pFilename)
}

func (p *sClient) GetRemoteFileProc(pCtx context.Context, pFriend string, pFilename string, pPersonal bool) (dto.IDownloadProcess, error) {
	return p.fRequester.GetRemoteFileProc(pCtx, pFriend, pFilename, pPersonal)
}

func (p *sClient) DelRemoteFileProc(pCtx context.Context, pFriend string, pFilename string, pPersonal bool) error {
	return p.fRequester.DelRemoteFileProc(pCtx, pFriend, pFilename, pPersonal)
}

func (p *sClient) GetRemoteListProc(pCtx context.Context) ([]dto.IDownloadProcess, error) {
	return p.fRequester.GetRemoteListProc(pCtx)
}
