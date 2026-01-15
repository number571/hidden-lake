package client

import (
	"context"
	"io"

	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/utils"
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

func (p *sClient) GetIndex(pCtx context.Context) (string, error) {
	return p.fRequester.GetIndex(pCtx)
}

func (p *sClient) GetFileInfo(pCtx context.Context, pAliasName string, pName string) (utils.IFileInfo, error) {
	return p.fRequester.GetFileInfo(pCtx, pAliasName, pName)
}

func (p *sClient) GetListFiles(pCtx context.Context, pAliasName string, pPage uint64) ([]utils.IFileInfo, error) {
	return p.fRequester.GetListFiles(pCtx, pAliasName, pPage)
}

func (p *sClient) DownloadFile(pW io.Writer, pCtx context.Context, pAliasName, pName string) error {
	return p.fRequester.DownloadFile(pW, pCtx, pAliasName, pName)
}
