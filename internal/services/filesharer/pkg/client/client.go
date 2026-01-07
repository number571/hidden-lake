package client

import (
	"context"
)

var (
	_ IClient = &sClient{}
)

type sClient struct {
	fBuilder   IBuilder
	fRequester IRequester
}

func NewClient(pBuilder IBuilder, pRequester IRequester) IClient {
	return &sClient{
		fBuilder:   pBuilder,
		fRequester: pRequester,
	}
}

func (p *sClient) GetFileInfo(pCtx context.Context, pAliasName string, pName string) (IFileInfo, error) {
	return p.fRequester.GetFileInfo(pCtx, pAliasName, p.fBuilder.GetFileInfo(pName))
}

func (p *sClient) GetListFiles(pCtx context.Context, pAliasName string, pPage uint64) ([]IFileInfo, error) {
	return p.fRequester.GetListFiles(pCtx, pAliasName, p.fBuilder.GetListFiles(pPage))
}

func (p *sClient) LoadFileChunk(pCtx context.Context, pAliasName, pName string, pChunk uint64) ([]byte, error) {
	return p.fRequester.LoadFileChunk(pCtx, pAliasName, p.fBuilder.LoadFileChunk(pName, pChunk))
}
