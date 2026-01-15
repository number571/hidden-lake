package client

import (
	"context"
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

func (p *sClient) SendPing(pCtx context.Context, pAliasName string) error {
	return p.fRequester.SendPing(pCtx, pAliasName)
}
