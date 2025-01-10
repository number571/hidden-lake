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

func (p *sClient) Ping(pCtx context.Context, pAliasName string) error {
	return p.fRequester.Ping(pCtx, pAliasName, p.fBuilder.Ping())
}
