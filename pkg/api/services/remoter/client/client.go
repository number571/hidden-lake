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

func (p *sClient) GetIndex(pCtx context.Context) (string, error) {
	return p.fRequester.GetIndex(pCtx)
}

func (p *sClient) ExecCommand(pCtx context.Context, pFriend string, pCommand ...string) ([]byte, error) {
	return p.fRequester.ExecCommand(pCtx, pFriend, p.fBuilder.ExecCommand(pCommand...))
}
