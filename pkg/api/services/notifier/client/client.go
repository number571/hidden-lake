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

func (p *sClient) GetIndex(pCtx context.Context) error {
	return p.fRequester.GetIndex(pCtx)
}

func (p *sClient) GetNotifyLimit(pCtx context.Context) (uint64, error) {
	return p.fRequester.GetNotifyLimit(pCtx)
}

func (p *sClient) SendNotification(pCtx context.Context, pAliasName string, pBody []byte) error {
	return p.fRequester.SendNotification(pCtx, pAliasName, pBody)
}

func (p *sClient) RecvNotification(pCtx context.Context, pSubscribeID string) (string, []byte, error) {
	return p.fRequester.RecvNotification(pCtx, pSubscribeID)
}
