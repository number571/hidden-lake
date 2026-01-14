package client

import (
	"context"

	"github.com/number571/hidden-lake/internal/services/messenger/pkg/message"
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

func (p *sClient) PushMessage(pCtx context.Context, pAliasName string, pBody string) (string, error) {
	return p.fRequester.PushMessage(pCtx, pAliasName, pBody)
}

func (p *sClient) LoadMessages(pCtx context.Context, pAliasName string, pPage uint64, pOffset uint64) ([]message.IMessage, error) {
	return p.fRequester.LoadMessages(pCtx, pAliasName, pPage, pOffset)
}

func (p *sClient) ListenChat(pCtx context.Context, pSubscribeID, pAliasName string) (message.IMessage, error) {
	return p.fRequester.ListenChat(pCtx, pSubscribeID, pAliasName)
}
